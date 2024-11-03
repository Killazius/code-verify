package compilation

import (
	"bytes"
	"compile-server/internal/models"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/exec"
	"time"
)

func MakePY(taskName string, userFile string) (string, error) {
	pathTask := fmt.Sprintf("src/%v", taskName)
	baseFile := fmt.Sprintf("%v/%v", pathTask, models.BasePy)
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

func TestPY(TaskName string, outputFile string) (string, error) {
	path := fmt.Sprintf("src/%v/%v", TaskName, models.TestPy)

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
			return "", fmt.Errorf("%s", stderrBuf.String())
		}
	}
	return "", nil
}
func RunPY(conn *websocket.Conn, userFile string, TaskName string) error {
	outputFile, err := MakePY(TaskName, userFile)
	if err != nil && outputFile == "" {
		conn.WriteJSON(models.Answer{
			Stage:   "build",
			Message: err.Error(),
		})
		return err
	}
	output, errCmd := TestPY(TaskName, outputFile)
	if errCmd != nil {
		err = os.Remove(outputFile)
		if err != nil {
			return fmt.Errorf("%s: %v", outputFile, err)
		}
		return errCmd
	}
	conn.WriteJSON(models.Answer{
		Stage:   "test",
		Message: output,
	})
	return nil
}
