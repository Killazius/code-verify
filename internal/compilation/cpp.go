package compilation

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func MakeCPPfile(taskName string, userFile string) error {
	baseFile := fmt.Sprintf("src/%v/base.cpp", taskName)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return err
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", userFile, err)
		return err
	}
	err = os.Remove(userFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, append(baseContent, userContent...), 0644)
	if err != nil {
		log.Printf("Ошибка записи в файл %s: %v\n", userFile, err)
		return err
	}

	userFileExe, err := CompileCPPfile(userFile, taskName)
	if err != nil {
		log.Printf("Файл не скомпилирован.")
		return err
	}
	err = TestCPPfile(userFileExe, taskName)
	if err != nil {
		log.Printf("Ошибка во время тестирования: %v", err)
		return err
	}
	outputFileExePath := fmt.Sprintf("src/%v/%v", taskName, userFileExe)
	err = os.Remove(outputFileExePath)
	if err != nil {
		return err
	}
	return nil
}

func CompileCPPfile(userFile string, TaskName string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	path := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	cmd := exec.Command("g++", "-o", path, userFile)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	err = os.Remove(userFile)
	if err != nil {
		return "", err
	}
	return userFileExe, nil

}

func TestCPPfile(userFile string, TaskName string) error {
	path := fmt.Sprintf("src/%v/test_cpp.go", TaskName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, userFile)
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
