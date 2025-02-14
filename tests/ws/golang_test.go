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

type testCaseGo struct {
	name       string
	message    ws.UserMessage
	compileErr bool
	testErr    bool
}

func TestWsGo(t *testing.T) {
	u := fmt.Sprintf("ws://%v/ws", host)
	tests := []testCaseGo{
		{
			name: "correct decision",
			message: ws.UserMessage{
				Code:   "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a)\n\tfmt.Scanln(&b)\n\tfmt.Println(a + b)\n}",
				Lang:   compilation.LangGo,
				TaskID: "1",
				Token:  "GO1",
			},
			compileErr: false,
			testErr:    false,
		},
		{
			name: "syntax error",
			message: ws.UserMessage{
				Code:   "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a)\n\tfmt.Scanln(&b)\n\tfmt.Println(ab)\n}",
				Lang:   compilation.LangGo,
				TaskID: "1",
				Token:  "GO2",
			},
			compileErr: true,
		},
		{
			name: "incorrect decision",
			message: ws.UserMessage{
				Code:   "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a)\n\tfmt.Scanln(&b)\n\tfmt.Println(a)\n}",
				Lang:   compilation.LangGo,
				TaskID: "1",
				Token:  "GO3",
			},
			compileErr: false,
			testErr:    true,
		},
		{
			name: "endless loop",
			message: ws.UserMessage{
				Code:   "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a)\n\tfmt.Scanln(&b)\n\tfor {\n\t\tfmt.Println(a+b)\n\t}\n}",
				Lang:   compilation.LangGo,
				TaskID: "1",
				Token:  "GO4",
			},
			compileErr: false,
			testErr:    true,
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
				if response.Stage == utils.Compile {
					if tt.compileErr {
						assert.NotEqual(t, "OK", response.Message)
						break
					}
					require.Equal(t, "OK", response.Message)

				} else if response.Stage == utils.Test {
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
