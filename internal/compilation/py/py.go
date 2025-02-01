package py

import (
	"compile-server/internal/compilation/test"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func Run(conn *websocket.Conn, userFile, taskName string) (*utils.CompilationResult, error) {
	const op = "compilation.py.Run"

	command := "py"
	output, errCmd := test.Run(command, taskName, userFile)
	if errCmd != nil {
		errSend := utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return nil, fmt.Errorf("%s: %w", op, errSend)
		}
		return nil, fmt.Errorf("%s: %w", op, errCmd)
	}
	errSend := utils.SendJSON(conn, utils.Test, output)
	if errSend != nil {
		return nil, fmt.Errorf("%s: %w", op, errSend)
	}
	err := os.Remove(userFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &utils.CompilationResult{Success: output == utils.OK, Output: output}, nil
}
