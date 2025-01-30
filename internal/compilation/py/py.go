package py

import (
	"compile-server/internal/compilation/test"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func Run(conn *websocket.Conn, userFile string) error {
	const op = "compilation.py.Run"

	command := "py"
	output, errCmd := test.Run(command, userFile)
	if errCmd != nil {
		errSend := utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, errCmd)
	}
	errSend := utils.SendJSON(conn, utils.Test, output)
	if errSend != nil {
		return fmt.Errorf("%s: %v", op, errSend)
	}
	err := os.Remove(userFile)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
