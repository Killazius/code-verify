package cpp

import (
	"compile-server/internal/compilation/test"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/exec"
	"strings"
)

func Compile(userFile string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	cmd := exec.Command("g++", "-o", userFileExe, userFile)

	output, errCmd := cmd.CombinedOutput()
	if errCmd != nil {
		removeErr := os.Remove(userFile)
		if removeErr != nil {
			return "", removeErr
		}
		return "", fmt.Errorf("%s", output)
	}
	err := os.Remove(userFile)
	if err != nil {
		return "", err
	}
	return userFileExe, nil

}

func Run(conn *websocket.Conn, userFile string) error {
	const op = "compilation.cpp.Run"

	userFileExe, err := Compile(userFile)

	if err != nil && userFileExe == "" {
		errSend := utils.SendJSON(conn, utils.Compile, err.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, err)
	}

	errSend := utils.SendJSON(conn, utils.Compile, utils.Success)
	if errSend != nil {
		return fmt.Errorf("%s: %v", op, errSend)
	}

	command := fmt.Sprintf("./%v", userFileExe)
	output, errCmd := test.Run(command)

	if errCmd != nil {
		errSend = utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, errCmd)
	}
	errSend = utils.SendJSON(conn, utils.Test, output)
	if errSend != nil {
		return fmt.Errorf("%s: %v", op, errSend)
	}

	err = os.Remove(userFileExe)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
