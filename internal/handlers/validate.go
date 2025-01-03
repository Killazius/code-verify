package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	Username string `json:"username"`
}

func GetName(log *slog.Logger, token string) (string, int) {
	url := fmt.Sprintf("https://studyingit-api.ru/api/code/auth/%s/", token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		log.Error("get request failed:", err)
		return "", http.StatusBadRequest
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Error("status code != 200:", resp.StatusCode)
		return "", http.StatusBadRequest
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("status code != 200:", resp.StatusCode)
		return "", http.StatusInternalServerError
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", http.StatusInternalServerError
	}

	return response.Username, http.StatusOK
}
