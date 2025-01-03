package utils

import "github.com/gorilla/websocket"

func SendStatus(conn *websocket.Conn, status int) error {
	message := struct {
		Status int `json:"status"`
	}{
		Status: status,
	}
	return conn.WriteJSON(message)
}

type Stage string

const (
	Build   Stage  = "build"
	Compile Stage  = "compile"
	Test    Stage  = "test"
	Success string = "success"
	Timeout string = "timeout"
)

func SendJSON(conn *websocket.Conn, stage Stage, message interface{}) error {
	json := struct {
		Stage   Stage  `json:"stage"`
		Message string `json:"message"`
	}{
		Stage:   stage,
		Message: message.(string),
	}

	return conn.WriteJSON(json)

}
