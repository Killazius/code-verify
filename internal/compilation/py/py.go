package py

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

func Build(taskName string, userFile string) (string, error) {
	pathTask := fmt.Sprintf("src/%v", taskName)
	baseFile := fmt.Sprintf("%v/%v", pathTask, compilation.BasePy)
	outputFile := fmt.Sprintf("%v/%v", pathTask, userFile)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		return "", fmt.Errorf("%s: %v", baseFile, err)
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		return "", fmt.Errorf("%s: %v", userFile, err)
	}

	err = os.WriteFile(outputFile, append(userContent, baseContent...), 0644)
	if err != nil {
		return "", fmt.Errorf("%s: %v", outputFile, err)
	}

	err = os.Remove(userFile)
	if err != nil {
		return "", fmt.Errorf("%s: %v", userFile, err)
	}

	return outputFile, nil
}

func Test(TaskName string, outputFile string) (string, error) {
	path := fmt.Sprintf("src/%v/%v", TaskName, compilation.TestPy)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, outputFile)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

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
	const op = "compilation.py.Run"
	outputFile, err := Build(TaskName, userFile)
	if err != nil && outputFile == "" {
		errSend := utils.SendJSON(conn, utils.Build, err.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, err)
	} else {
		errSend := utils.SendJSON(conn, utils.Build, utils.Success)
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
	}
	output, errCmd := Test(TaskName, outputFile)
	output = strings.ReplaceAll(output, "\n", "")
	if errCmd != nil {
		err = os.Remove(outputFile)
		if err != nil {
			return fmt.Errorf("%s: %v", outputFile, err)
		}
		errSend := utils.SendJSON(conn, utils.Test, errCmd.Error())
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
		return fmt.Errorf("%s: %v", op, errCmd)
	} else {
		errSend := utils.SendJSON(conn, utils.Test, output)
		if errSend != nil {
			return fmt.Errorf("%s: %v", op, errSend)
		}
	}
	return nil
}
