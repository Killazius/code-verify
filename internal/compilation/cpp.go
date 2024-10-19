package compilation

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func MakeCPPfile(taskName string, userFile string) {
	baseFile := fmt.Sprintf("src/%v/base.cpp", taskName)
	outputFile := "solution.cpp"
	outputFileExe := "solution.exe"

	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return
	}

	sumContent, err := os.ReadFile(userFile)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v\n", userFile, err)
		return
	}

	err = os.WriteFile(outputFile, append(baseContent, sumContent...), 0644)
	if err != nil {
		log.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
		return
	}

	err = CompileCPPfile(outputFileExe, outputFile, taskName)
	if err != nil {
		log.Printf("Файл не скомпилирован.")
		return
	}
	err = TestCPPfile(taskName)
	outputFileExePath := fmt.Sprintf("src/%v/%v", taskName, outputFileExe)
	os.Remove(outputFileExePath)
	if err != nil {
		log.Printf("Ошибка во время тестирования: %v", err)
		return
	}

}

func CompileCPPfile(outputFileExe string, outputFile string, TaskName string) error {
	path := fmt.Sprintf("src/%v/%v", TaskName, outputFileExe)
	cmd := exec.Command("g++", "-o", path, outputFile)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Remove("user.cpp")
	if err != nil {
		return err
	}
	err = os.Remove("solution.cpp")
	if err != nil {
		return err
	}
	return nil

}

func TestCPPfile(TaskName string) error {
	path := fmt.Sprintf("src/%v/test.go", TaskName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", path)
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
