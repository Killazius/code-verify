package utils

import "github.com/gorilla/websocket"

type StatusMessage struct {
	Status int `json:"status"`
}

func SendStatus(conn *websocket.Conn, status int) error {
	message := StatusMessage{
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

type StageMessage struct {
	Stage   Stage  `json:"stage"`
	Message string `json:"message"`
}

func SendJSON(conn *websocket.Conn, stage Stage, message interface{}) error {
	json := StageMessage{
		Stage:   stage,
		Message: message.(string),
	}
	return conn.WriteJSON(json)
}
