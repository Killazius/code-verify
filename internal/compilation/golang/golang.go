package golang

import (
	"compile-server/internal/compilation/test"
	"compile-server/internal/handlers/ws/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/exec"
	"strings"
)

func compile(userFile string) (string, error) {
	userFileExe := strings.Replace(userFile, ".go", ".exe", 1)
	cmd := exec.Command("go", "build", userFile)
	output, err := cmd.CombinedOutput()
	defer func() {
		err := os.Remove(userFile)
		if err != nil {
			return
		}
	}()
	if err != nil {
		return "", fmt.Errorf("%s", output)
	}
	return userFileExe, nil
}

func CompileAndRun(conn *websocket.Conn, userFile, taskName string) (*utils.CompilationResult, error) {
	const op = "compilation.cpp.Run"

	userFileExe, err := compile(userFile)
	defer func() {
		err := os.Remove(userFileExe)
		if err != nil {
			return
		}
	}()

	if err != nil && userFileExe == "" {
		errSend := utils.SendJSON(conn, utils.Compile, err.Error())
		if errSend != nil {
			return nil, fmt.Errorf("%s: %w", op, errSend)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	errSend := utils.SendJSON(conn, utils.Compile, utils.OK)
	if errSend != nil {
		return nil, fmt.Errorf("%s: %w", op, errSend)
	}

	command := fmt.Sprintf("./%v", userFileExe)
	output, errCmd := test.Run(command, taskName)

	if errCmd != nil {
		errSend = utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return nil, fmt.Errorf("%s: %w", op, errSend)
		}
		// TODO здесь не всегда timeout может быть. надо бы додумать и переделать
		return &utils.CompilationResult{Success: output == utils.OK, Output: utils.Timeout}, fmt.Errorf("%s: %w", op, errCmd)
	}
	errSend = utils.SendJSON(conn, utils.Test, output)
	if errSend != nil {
		return nil, fmt.Errorf("%s: %w", op, errSend)
	}

	return &utils.CompilationResult{Success: output == utils.OK, Output: output}, nil
}
