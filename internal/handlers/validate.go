package handlers

import (
	"bytes"
	"compile-server/internal/config"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}
func GetID(token string) (string, int, error) {
	const op = "handlers.validate.GetID"
	if config.Env == config.Local {
		return token, http.StatusOK, nil
	}
	client := newClient()

	url := fmt.Sprintf("https://studyingit-api.ru/api/code/auth/%s/", token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("%s: %w", op, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", http.StatusBadRequest, fmt.Errorf("%s: %v", op, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	var response struct {
		UserID string `json:"user_id"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	return response.UserID, http.StatusOK, nil
}

type reqBody struct {
	UserID string `json:"user_id"`
	TaskID string `json:"task_id"`
}

func MarkTaskAsCompleted(userID, taskID string) (int, error) {
	const op = "handlers.validate.MarkTaskAsCompleted"
	if config.Env == config.Local {
		return http.StatusOK, nil
	}
	url := "https://studyingit-api.ru/api/complete/"

	client := newClient()

	body, err := json.Marshal(reqBody{UserID: userID, TaskID: taskID})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("%s: %v", op, resp.Status)
	}
	return resp.StatusCode, nil
}
