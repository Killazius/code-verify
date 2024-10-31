package handlers

import (
	"compile-server/internal/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Received: %s", message)
		var userCode models.Code
		if err := json.Unmarshal(message, &userCode); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}
		log.Printf("Received code: %s", userCode.Code)
		log.Printf("Language: %s", userCode.Lang)
		log.Printf("Task Name: %s", userCode.TaskName)
		log.Printf("Username: %s", userCode.UserName)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			break
		}

	}
}
