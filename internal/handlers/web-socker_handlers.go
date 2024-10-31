package handlers

import (
	"compile-server/internal/compilation"
	"compile-server/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
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
		_, message, err := conn.ReadMessage()
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

		userFile := fmt.Sprintf("%v-%v.%v", userCode.TaskName, userCode.UserName, userCode.Lang)
		file, err := os.Create(userFile)
		if err != nil {
			log.Println("Error creating file:", err)
			continue
		}
		defer file.Close()

		_, err = file.WriteString(userCode.Code)
		if err != nil {
			log.Println("Error writing to file:", err)
			continue
		}

		switch userCode.Lang {
		case "cpp":
			{
				err := compilation.MakeCPPfile(conn, userCode.TaskName, userFile)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			}
		case "py":
			{
				err := compilation.MakePYfile(conn, userCode.TaskName, userFile)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			}
		default:
			http.Error(w, fmt.Sprint("Compilations have not started."), http.StatusBadRequest)
		}

	}
}
