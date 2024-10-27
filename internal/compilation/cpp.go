package compilation

import (
	"compile-server/internal/models"
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
		return models.HandleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", baseFile, err))
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", userFile, err))
	}
	err = os.Remove(userFile)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка при удалении временного файла: %v", err))
	}

	err = os.WriteFile(userFile, append(baseContent, userContent...), 0644)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка при записи в файл %s: %v", userFile, err))
	}

	userFileExe, err := CompileCPPfile(userFile, taskName)
	if err != nil || userFileExe == "" {
		return models.HandleCommonError(fmt.Errorf("файл не скомпилирован: %v", err))
	}
	err = TestCPPfile(userFileExe, taskName)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка во время тестирования: %v", err))
	}
	outputFileExePath := fmt.Sprintf("src/%v/%v", taskName, userFileExe)
	err = os.Remove(outputFileExePath)
	if err != nil {
		return models.HandleCommonError(fmt.Errorf("ошибка при удалении исполняемого файла: %v", err))
	}
	return nil
}

func CompileCPPfile(userFile string, TaskName string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	path := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	cmd := exec.Command("g++", "-o", path, userFile)

	err_cmd := cmd.Run()
	if err_cmd != nil {
		err := os.Remove(userFile)
		if err != nil {
			return "", err
		}
		return "", err_cmd
	}
	err := os.Remove(userFile)
	if err != nil {
		return "", err
	}
	return userFileExe, nil

}

func TestCPPfile(userFile string, TaskName string) error {
	path := fmt.Sprintf("src/%v/test_cpp.go", TaskName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path, userFile)
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
