//go:build integration

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

type testCaseCpp struct {
	name       string
	message    ws.UserMessage
	compileErr bool
	testErr    bool
}

func TestWsCpp(t *testing.T) {
	u := fmt.Sprintf("ws://%v/ws", host)
	tests := []testCaseCpp{
		{
			name: "correct decision",
			message: ws.UserMessage{
				Code:   "#include <iostream>\nint main() {int a,b; std::cin>>a>>b; std::cout<<a+b;}",
				Lang:   compilation.LangCpp,
				TaskID: "1",
				Token:  "CPP1",
			},
			compileErr: false,
			testErr:    false,
		},
		{
			name: "syntax error",
			message: ws.UserMessage{
				Code:   "#include <iostream>\nint main() {int a,b; std::cin>>a>>b; std::cout<<ab;}",
				Lang:   compilation.LangCpp,
				TaskID: "1",
				Token:  "CPP2",
			},
			compileErr: true,
		},
		{
			name: "incorrect decision",
			message: ws.UserMessage{
				Code:   "#include <iostream>\nint main() {int a,b; std::cin>>a>>b; std::cout<<a;}",
				Lang:   compilation.LangCpp,
				TaskID: "1",
				Token:  "CPP3",
			},
			compileErr: false,
			testErr:    true,
		},
		{
			name: "endless loop",
			message: ws.UserMessage{
				Code:   "#include <iostream>\nint main() {int a,b; std::cin>>a>>b; while (true);}",
				Lang:   compilation.LangCpp,
				TaskID: "1",
				Token:  "CPP4",
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
