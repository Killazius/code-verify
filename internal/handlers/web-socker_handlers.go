package handlers

import (
	"compile-server/internal/compilation"
	"compile-server/internal/models"
	"encoding/json"
	"fmt"
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
	log.Println("New connection from", r.RemoteAddr)
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("connection close error:", err)
		}
	}(conn)
	defer log.Printf("Connection closed from %s", r.RemoteAddr)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var user models.UserJson
		if err := json.Unmarshal(message, &user); err != nil {
			log.Println("Error decoding JSON:", err)
			break
		}

		userName, status := compilation.GetName(user.Token)
		err = conn.WriteJSON(models.TokenAnswer{
			Status: status},
		)
		if err != nil {
			return
		}
		if status != http.StatusOK {
			log.Println("token is invalid")
			break
		}

		userFile := fmt.Sprintf("%v-%v.%v", user.TaskName, userName, user.Lang)
		err = compilation.CreateFile(userFile, user.Code, user.Lang)
		if err != nil {
			log.Println(err)
		}

		switch user.Lang {
		case "cpp":
			{
				err = compilation.RunCPP(conn, userFile, user.TaskName)
				if err != nil {
					break
				}
			}
		case "py":
			{
				err = compilation.RunPY(conn, userFile, user.TaskName)
				if err != nil {
					break
				}
			}
		default:
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unsupported language: %s", user.Lang)))
			if err != nil {
				return
			}
		}
		break
	}
}
