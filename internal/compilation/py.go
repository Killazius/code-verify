package compilation

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func MakePYfile(taskName string, userFile string) {
	pathTask := fmt.Sprintf("src/%s", taskName)
	baseFile := fmt.Sprintf("%s/base.py", pathTask)
	outputFile := fmt.Sprintf("%s/%s", pathTask, userFile)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", userFile, err)
		return
	}

	err = os.WriteFile(outputFile, append(userContent, baseContent...), 0644)
	if err != nil {
		log.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
		return
	}

	err = os.Remove(userFile)
	if err != nil {
		log.Printf("Ошибка в удалении файла %s: %v\n", userFile, err)
	}

	err = TestPYfile(taskName, outputFile)
	outputFilePath := fmt.Sprint(outputFile)
	os.Remove(outputFilePath)
	if err != nil {
		log.Printf("Ошибка во время тестирования: %v", err)
		return
	}
}

func TestPYfile(TaskName string, outputFile string) error {
	path := fmt.Sprintf("src/%v/test_py.go", TaskName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("command finished with error: %v", err)
			return fmt.Errorf("command finished with error: %w", err)
		}
		log.Println("command finished successfully")
	case <-ctx.Done():

		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		log.Println("command timed out")
		return fmt.Errorf("command timed out")
	}

	return nil
}
