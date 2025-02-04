package ws_test

import (
	"compile-server/internal/compilation"
	"compile-server/internal/handlers/ws"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	host = "localhost:1234"
)

type testCasePy struct {
	name    string
	message ws.UserMessage
	testErr bool
}

func TestWebSocketPY(t *testing.T) {
	u := fmt.Sprintf("ws://%v/ws", host)
	tests := []testCasePy{
		{
			name: "correct decision on py",
			message: ws.UserMessage{
				Code:     "a=int(input())\nb=int(input())\nprint(a+b)",
				Lang:     compilation.LangPy,
				TaskName: "1-sum",
				Token:    "PY1",
			},
			testErr: false,
		},
		{
			name: "syntax error on py",
			message: ws.UserMessage{
				Code:     "a=int(input())\nb=int(input())\nprint(ab)",
				Lang:     compilation.LangPy,
				TaskName: "1-sum",
				Token:    "PY2",
			},
			testErr: true,
		},
		{
			name: "incorrect decision on py",
			message: ws.UserMessage{
				Code:     "a=int(input())\nb=int(input())\nprint(a)",
				Lang:     compilation.LangPy,
				TaskName: "1-sum",
				Token:    "PY3",
			},
			testErr: true,
		},
		{
			name: "endless loop on py",
			message: ws.UserMessage{
				Code:     "a = int(input())\nb = int(input())\nwhile True:\n\tprint(a)",
				Lang:     compilation.LangPy,
				TaskName: "1-sum",
				Token:    "PY4",
			},
			testErr: true,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			webSocket, _, err := websocket.DefaultDialer.Dial(u, nil)
			require.NoError(t, err)
			defer webSocket.Close()

			err = webSocket.WriteJSON(tt.message)
			require.NoError(t, err)

			var MessageStatus utils.StatusMessage
			err = webSocket.ReadJSON(&MessageStatus)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, MessageStatus.Status)

			var response utils.StageMessage
			for {
				err = webSocket.ReadJSON(&response)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					break
				}
				t.Log(response)
				require.NoError(t, err)

				if response.Stage == utils.Test {
					if tt.testErr {
						assert.NotEqual(t, "OK", response.Message)
						break
					}
					assert.Equal(t, "OK", response.Message)

				}
				if response.Stage == utils.Test {
					break
				}
			}
		})
	}
}
