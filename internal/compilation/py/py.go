package py

import (
	"compile-server/internal/compilation/test"
	"compile-server/internal/config"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func Run(conn *websocket.Conn, userFile, taskName string) (*utils.CompilationResult, error) {
	const op = "compilation.py.Run"
	defer func() {
		err := os.Remove(userFile)
		if err != nil {
			return
		}
	}()
	var command string
	if config.Env == config.Local {
		command = "py"
	} else {
		command = "python3"
	}
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
	return &utils.CompilationResult{Success: output == utils.OK, Output: output}, nil
}
