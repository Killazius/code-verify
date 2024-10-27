package compilation

import (
	"compile-server/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func MakePYfile(taskName string, userFile string) error {
	pathTask := fmt.Sprintf("src/%s", taskName)
	baseFile := fmt.Sprintf("%s/base.py", pathTask)
	outputFile := fmt.Sprintf("%s/%s", pathTask, userFile)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", baseFile, err))
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", userFile, err))
	}

	err = os.WriteFile(outputFile, append(userContent, baseContent...), 0644)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", outputFile, err))
	}

	err = os.Remove(userFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка в удалении файла %s: %v", userFile, err))
	}

	err = TestPYfile(taskName, outputFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка во время тестирования: %v", err))
	}
	outputFilePath := fmt.Sprint(outputFile)
	err = os.Remove(outputFilePath)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка в удалении файла %s: %v", outputFilePath, err))
	}
	return nil
}

func TestPYfile(TaskName string, outputFile string) error {
	path := fmt.Sprintf("src/%v/test_py.go", TaskName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
		log.Println("tests passed")
	case <-ctx.Done():

		if err := cmd.Process.Kill(); err != nil {
			return err
		}
		log.Println("tests failed")
	}

	return nil
}
