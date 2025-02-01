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
	Compile Stage  = "compile"
	Test    Stage  = "test"
	OK      string = "OK"
	Timeout string = "timeout"
)

type CompilationResult struct {
	Success bool
	Output  string
}

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
