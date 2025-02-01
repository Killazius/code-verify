package handlers

import (
	"compile-server/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetName(token string) (string, int, error) {
	const op = "handlers.validate.GetName"
	if config.Env == "local" {
		return "localhost", http.StatusOK, nil
	}
	url := fmt.Sprintf("https://studyingit-api.ru/api/code/auth/%s/", token)
	client := &http.Client{}
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
		Username string `json:"username"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	return response.Username, http.StatusOK, nil
}

func MarkTaskAsCompleted(username, token string) (int, error) {
	const op = "handlers.validate.MarkTaskAsCompleted"
	if config.Env == "local" {
		return http.StatusOK, nil
	}
	url := fmt.Sprintf("https://studyingit-api.ru/api/%v/complete/", username)
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("%s: %v", op, resp.Status)
	}
	return resp.StatusCode, nil
}
