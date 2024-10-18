package compilation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	cmd := exec.Command("go", "run", path)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}
	log.Println(string(output))
	return nil
}
