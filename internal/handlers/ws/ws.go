package ws

import (
	"compile-server/internal/compilation"
	"compile-server/internal/logger"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"log/slog"
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

type UserMessage struct {
	Code     string `json:"code"`
	Lang     string `json:"lang"`
	TaskName string `json:"task_name"`
	Token    string `json:"token"`
}

// New TODO: Решить проблему с логированием.
// Почему то если сделать 2 запроса, то во 2 логе кидает op и request_id первого запроса
// upd: вроде решил, но с помощью затенения переменной log
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ws.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("fail to upgrade connection", slog.Any(logger.Err, err))
		}

		defer func(conn *websocket.Conn) {
			if err := conn.Close(); err != nil {
				log.Error("connection close", slog.Any(logger.Err, err))
				return
			}
		}(conn)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("could not read the message", slog.Any(logger.Err, err))
				break
			}

			var user UserMessage
			if err = json.Unmarshal(message, &user); err != nil {
				log.Error("unmarshal failed", slog.Any(logger.Err, err))
				break
			}
			var (
				userName string
				status   int
			)
			// TODO: mock for token. Надо будет рефакторить, не нравится решение
			if os.Getenv("ENV") == "local" {
				userName = "localhost"
				status = http.StatusOK
			} else {
				userName, status = compilation.GetName(user.Token)
			}
			// как то не читаемо, может убрать???
			err = conn.WriteJSON(struct {
				Status int `json:"status"`
			}{Status: status})
			if err != nil {
				log.Error("send status-json failed", slog.Any(logger.Err, err))
				return
			}
			if status != http.StatusOK {
				log.Error("token-status != 200", slog.Any(logger.Err, err))
				break
			}

			userFile := fmt.Sprintf("%v-%v.%v", user.TaskName, userName, user.Lang)
			err = compilation.CreateFile(userFile, user.Code, user.Lang)
			if err != nil {
				log.Error("create file failed", slog.Any(logger.Err, err))
			}

			switch user.Lang {
			case "cpp":
				{
					err = compilation.RunCPP(conn, userFile, user.TaskName)
					if err != nil {
						log.Error("run cpp file failed", slog.Any(logger.Err, err))
						break
					}
				}
			case "py":
				{
					err = compilation.RunPY(conn, userFile, user.TaskName)
					if err != nil {
						log.Error("run py file failed", slog.Any(logger.Err, err))
						break
					}
				}
			default:
				err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unsupported language: %s", user.Lang)))
				if err != nil {
					log.Error("failed writeMessage to conn", slog.Any(logger.Err, err))
					return
				}
			}
			break
		}
	}
}
