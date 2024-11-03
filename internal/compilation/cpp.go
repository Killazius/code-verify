package compilation

import (
	"bytes"
	"compile-server/internal/models"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/exec"
	"strings"
	"time"
)

func BuildCPP(taskName string, userFile string) error {
	baseFile := fmt.Sprintf("src/%v/%v", taskName, models.BaseCpp)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		return fmt.Errorf("%s: %v", baseFile, err)
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		return fmt.Errorf("%s: %v", userFile, err)
	}
	err = os.Remove(userFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, append(baseContent, userContent...), 0644)
	if err != nil {
		return fmt.Errorf("%s: %v", userFile, err)
	}
	return nil
}

func CompileCPP(userFile string, TaskName string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	path := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	cmd := exec.Command("g++", "-o", path, userFile)

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

func TestCPP(userFile string, TaskName string) (string, error) {
	path := fmt.Sprintf("src/%v/%v", TaskName, models.TestCpp)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, userFile)
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
			return fmt.Sprintf("Error: %vnOutput: %snErrors: %s", err, stdoutBuf.String(), stderrBuf.String()), nil
		}
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			return "", err
		}
	}

	return "", nil
}

func RunCPP(conn *websocket.Conn, userFile string, TaskName string) error {
	err := BuildCPP(TaskName, userFile)
	if err != nil {
		conn.WriteJSON(models.Answer{
			Stage:   "build",
			Message: err.Error(),
		})
		return err
	}
	userFileExe, err := CompileCPP(userFile, TaskName)
	if err != nil && userFileExe == "" {
		conn.WriteJSON(models.Answer{
			Stage:   "compile",
			Message: err.Error(),
		})
		return err
	}
	output, err := TestCPP(userFileExe, TaskName)
	if err != nil {
		return err
	}
	outputFileExePath := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	err = os.Remove(outputFileExePath)
	if err != nil {
		return err
	}
	conn.WriteJSON(models.Answer{
		Stage:   "test",
		Message: output,
	})
	return nil
}
