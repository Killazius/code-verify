package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetName(token string, env string) (string, int, error) {
	const op = "handlers.validate.GetName"
	if env == "local" {
		return "localhost", http.StatusOK, nil
	}
	url := fmt.Sprintf("https://studyingit-api.ru/api/code/auth/%s/", token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("%s: %v", op, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", http.StatusBadRequest, fmt.Errorf("%s: %s", op, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %v", op, err)
	}

	var response struct {
		Username string `json:"username"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s: %v", op, err)
	}

	return response.Username, http.StatusOK, nil
}
