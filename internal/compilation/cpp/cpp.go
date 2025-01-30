package cpp

import (
	"bytes"
	"compile-server/internal/compilation"
	"compile-server/internal/handlers/ws/utils"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/exec"
	"strings"
	"time"
)

//func Build(taskName string, userFile string) error {
//	//baseFile := fmt.Sprintf("src/%v/%v", taskName, compilation.BaseCpp)
//	//
//	//baseContent, err := os.ReadFile(baseFile)
//	//if err != nil {
//	//	return fmt.Errorf("%s: %v", baseFile, err)
//	//}
//
//	userContent, err := os.ReadFile(userFile)
//	if err != nil {
//		return fmt.Errorf("%s: %v", userFile, err)
//	}
//	err = os.Remove(userFile)
//	if err != nil {
//		return err
//	}
//
//	err = os.WriteFile(userFile, userContent, 0600)
//	if err != nil {
//		return fmt.Errorf("%s: %v", userFile, err)
//	}
//	return nil
//}

func Compile(userFile string, TaskName string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	//path := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
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

func Test(userFile string, TaskName string) (string, error) {
	path := fmt.Sprintf("src/%v/%v", TaskName, compilation.TestCpp)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, userFile)
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	if err := cmd.Start(); err != nil {
		return "", err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", err
		}
		return stdoutBuf.String(), nil
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			return "", err
		}
		return utils.Timeout, nil
	}
}

func Run(conn *websocket.Conn, userFile string, TaskName string) error {
	const op = "compilation.cpp.Run"
	//if err := Build(TaskName, userFile); err != nil {
	//	errSend := utils.SendJSON(conn, utils.Build, err.Error())
	//	if errSend != nil {
	//		return fmt.Errorf("%s: %v", op, errSend)
	//	}
	//	return fmt.Errorf("%s: %v", op, err)
	//}
	//errSend := utils.SendJSON(conn, utils.Build, utils.Success)
	//if errSend != nil {
	//	return fmt.Errorf("%s: %v", op, errSend)
	//}

	userFileExe, err := Compile(userFile, TaskName)

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

	output, errCmd := Test(userFileExe, TaskName)
	output = strings.ReplaceAll(output, "\n", "")

	defer func() {
		outputFileExePath := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
		err = os.Remove(outputFileExePath)
		if err != nil {
			return
		}
	}()

	if errCmd != nil {
		errSend := utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, errCmd)
	}
	errSend = utils.SendJSON(conn, utils.Test, output)
	if errSend != nil {
		return fmt.Errorf("%s: %v", op, errSend)
	}

	outputFileExePath := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	err = os.Remove(outputFileExePath)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
