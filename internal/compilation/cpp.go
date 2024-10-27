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

type CompilationError struct {
	Msg    string
	Reason error
}

func (e *CompilationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Reason)
}

func handleCommonError(err error) error {
	if err != nil {
		return &CompilationError{
			Msg:    "Во время компиляции произошла ошибка.",
			Reason: err,
		}
	}
	return nil
}

func MakeCPPfile(taskName string, userFile string) error {
	baseFile := fmt.Sprintf("src/%v/base.cpp", taskName)

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", baseFile, err))
	}

	userContent, err := os.ReadFile(userFile)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка чтения файла %s: %v", userFile, err))
	}
	err = os.Remove(userFile)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка при удалении временного файла: %v", err))
	}

	err = os.WriteFile(userFile, append(baseContent, userContent...), 0644)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка при записи в файл %s: %v", userFile, err))
	}

	userFileExe, err := CompileCPPfile(userFile, taskName)
	if err != nil {
		return handleCommonError(fmt.Errorf("файл не скомпилирован: %v", err))
	}
	err = TestCPPfile(userFileExe, taskName)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка во время тестирования: %v", err))
	}
	outputFileExePath := fmt.Sprintf("src/%v/%v", taskName, userFileExe)
	err = os.Remove(outputFileExePath)
	if err != nil {
		return handleCommonError(fmt.Errorf("ошибка при удалении исполняемого файла: %v", err))
	}
	return nil
}

func CompileCPPfile(userFile string, TaskName string) (string, error) {
	userFileExe := strings.Replace(userFile, ".cpp", ".exe", 1)
	path := fmt.Sprintf("src/%v/%v", TaskName, userFileExe)
	cmd := exec.Command("g++", "-o", path, userFile)

	err := cmd.Run()
	if err != nil {
		return "", handleCommonError(fmt.Errorf("ошибка компиляции: %v", err))
	}
	err = os.Remove(userFile)
	if err != nil {
		return "", handleCommonError(fmt.Errorf("ошибка при удалении исходного файла: %v", err))
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
		return handleCommonError(fmt.Errorf("ошибка при запуске тестов: %v", err))
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return handleCommonError(fmt.Errorf("тест закончился с ошибкой: %v", err))
		}
		log.Println("Задача решена верно")
	case <-ctx.Done():

		if err := cmd.Process.Kill(); err != nil {
			return handleCommonError(fmt.Errorf("невозможно удалить процесс: %v", err))
		}
		return handleCommonError(fmt.Errorf("время на тестирование закончено. Задача не решена"))
	}

	return nil
}
