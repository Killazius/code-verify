package ws

import (
	"compile-server/internal/compilation"
	"compile-server/internal/compilation/cpp"
	"compile-server/internal/compilation/golang"
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
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ws.New"
		log := log.With(
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
			err := conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing connection"), time.Now().Add(time.Second))
			if err != nil {
				log.Error("failed to send close message", slog.String(logger.Err, err.Error()))
			}

			if err = conn.Close(); err != nil {
				log.Error("connection close", slog.String(logger.Err, err.Error()))
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
		log.Info("request JSON decoded",
			slog.String("code", user.Code),
			slog.String("lang", string(user.Lang)),
			slog.String("task_id", user.TaskID),
		)
		userID, status, errGet := handlers.GetID(user.Token)
		err = utils.SendStatus(conn, status)
		if err != nil {
			log.Error("send status-json failed", slog.String(logger.Err, err.Error()))
			return
		}
		if errGet != nil {
			log.Error("get userID failed", slog.String(logger.Err, errGet.Error()))
			return
		}

		if status != http.StatusOK {
			log.Error("token-status != 200", slog.Int("token-status", status), slog.String(logger.Err, err.Error()))
			return
		}
		log.Info("token verified", slog.String("userID", userID))

		userFile := fmt.Sprintf("%v-%v.%v", user.TaskID, userID, user.Lang)
		err = compilation.CreateFile(userFile, user.Code, user.Lang)
		if err != nil {
			log.Error("create file failed", slog.String(logger.Err, err.Error()))
		}

		result, err := handleLanguage(conn, userFile, user.Lang, user.TaskID)
		if err != nil {
			log.Error("handle language failed", slog.String(logger.Err, err.Error()))
		}
		log.Info("result is received", slog.Any("result", result))

		if result != nil && result.Success {
			status, err = handlers.MarkTaskAsCompleted(userID, user.TaskID)
			if err != nil || status != http.StatusOK {
				log.Info("mark task failed", slog.Int("token-status", status), slog.String(logger.Err, err.Error()))
				return
			}
			log.Info("task marked", slog.String("task", user.TaskID), slog.String("username", userID))
		}

	}
}

func handleLanguage(conn *websocket.Conn, userFile string, lang compilation.Lang, taskName string) (*utils.CompilationResult, error) {
	switch lang {
	case compilation.LangCpp:
		return cpp.CompileAndRun(conn, userFile, taskName)
	case compilation.LangPy:
		return py.Run(conn, userFile, taskName)
	case compilation.LangGo:
		return golang.CompileAndRun(conn, userFile, taskName)
	default:
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unsupported language: %s", lang)))
		if err != nil {
			return nil, fmt.Errorf("failed to write message to conn: %v", err)
		}
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}
