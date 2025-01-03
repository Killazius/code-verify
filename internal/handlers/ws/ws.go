package ws

import (
	"compile-server/internal/compilation"
	"compile-server/internal/compilation/cpp"
	"compile-server/internal/compilation/py"
	"compile-server/internal/handlers"
	"compile-server/internal/handlers/ws/utils"
	"compile-server/internal/logger"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UserMessage struct {
	Code     string           `json:"code"`
	Lang     compilation.Lang `json:"lang"`
	TaskName string           `json:"task_name"`
	Token    string           `json:"token"`
}

// New TODO: Решить проблему с логированием.
// Почему то если сделать 2 запроса, то во 2 логе кидает op и request_id первого запроса
// upd: вроде решил, но с помощью затенения переменной log
func New(log *slog.Logger, env string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ws.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("fail to upgrade connection", slog.String(logger.Err, err.Error()))
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err = conn.Close(); err != nil {
				log.Error("connection close", slog.String(logger.Err, err.Error()))
				return
			}
		}()

		for {
			var message []byte
			_, message, err = conn.ReadMessage()
			if err != nil {
				log.Error("could not read the message", slog.String(logger.Err, err.Error()))
				return
			}

			var user UserMessage
			if err = json.Unmarshal(message, &user); err != nil {
				log.Error("unmarshal failed", slog.String(logger.Err, err.Error()))
				return
			}
			log.Info("request JSON decoded", slog.Any("json", user))
			userName, status, errGet := handlers.GetName(user.Token, env)

			if errGet != nil {
				log.Error("get name failed", slog.String(logger.Err, errGet.Error()))
				return
			}
			err = utils.SendStatus(conn, status)
			if err != nil {
				log.Error("send status-json failed", slog.String(logger.Err, err.Error()))
				return
			}

			if status != http.StatusOK {
				log.Error("token-status != 200", slog.Int("token-status", status))
				return
			}
			log.Info("token verified", slog.String("username", userName))

			userFile := fmt.Sprintf("%v-%v.%v", user.TaskName, userName, user.Lang)
			err = compilation.CreateFile(userFile, user.Code, user.Lang)
			if err != nil {
				log.Error("create file failed", slog.String(logger.Err, err.Error()))
				return
			}

			switch user.Lang {
			case "cpp":
				{
					err = cpp.Run(conn, userFile, user.TaskName)
					if err != nil {
						log.Error("run cpp file failed", slog.String(logger.Err, err.Error()))
						return
					}
				}
			case "py":
				{
					err = py.Run(conn, userFile, user.TaskName)
					if err != nil {
						log.Error("run py file failed", slog.String(logger.Err, err.Error()))
						return
					}
				}
			default:
				err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unsupported language: %s", user.Lang)))
				if err != nil {
					log.Error("failed writeMessage to conn", slog.String(logger.Err, err.Error()))
					return
				}
			}
			break
		}
	}
}
