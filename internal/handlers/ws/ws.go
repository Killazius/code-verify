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
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type UserMessage struct {
	Code     string           `json:"code"`
	Lang     compilation.Lang `json:"lang"`
	TaskName string           `json:"task_name"`
	Token    string           `json:"token"`
}

func New(log *slog.Logger) http.HandlerFunc {
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
		userName, status, errGet := handlers.GetName(user.Token)

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
		}

		result, err := handleLanguage(conn, userFile, user.Lang)
		if err != nil {
			log.Error("handle language failed", slog.String(logger.Err, err.Error()))
		}
		log.Info("result is received", slog.Any("result", result))

		if result != nil && result.Success {
			_, err = handlers.MarkTaskAsCompleted(userName, user.Token)
			if err != nil {
				return
			}
		}

	}
}

func handleLanguage(conn *websocket.Conn, userFile string, lang compilation.Lang) (*utils.CompilationResult, error) {
	switch lang {
	case compilation.LangCpp:
		return cpp.CompileAndRun(conn, userFile)
	case compilation.LangPy:
		return py.Run(conn, userFile)
	default:
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unsupported language: %s", lang)))
		if err != nil {
			return nil, fmt.Errorf("failed to write message to conn: %v", err)
		}
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}
