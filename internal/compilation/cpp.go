package compilation

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func MakeCPPfile(taskName string, userFile string) {
	baseFile := fmt.Sprintf("src/%v/base.cpp", taskName)
	//sumFile := "sum_cpp.txt"
	outputFile := "solution.cpp"
	outputFileExe := "solution.exe"

	// Чтение содержимого base.cpp
	baseContent, err := ioutil.ReadFile(baseFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return
	}

	// Чтение содержимого user.cpp
	sumContent, err := ioutil.ReadFile(userFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", userFile, err)
		return
	}

	// Создание файла main.cpp
	err = ioutil.WriteFile(outputFile, append(baseContent, sumContent...), 0644)
	if err != nil {
		fmt.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
		return
	}
	//fmt.Printf("Файл %s успешно создан.\n", outputFile)

	err = CompileCPPfile(outputFileExe, outputFile, taskName)
	if err != nil {
		log.Println("Файл не скомпилирован.")
		return
	}
	err = TestCPPfile(taskName)
	if err != nil {
		log.Println("Ошибка во время тестирования.")
		return
	}

}

func CompileCPPfile(outputFileExe string, outputFile string, Task_Name string) error {
	path := fmt.Sprintf("src/%v/%v", Task_Name, outputFileExe)
	cmd := exec.Command("g++", "-o", path, outputFile)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	os.Remove("user.cpp")
	os.Remove("solution.cpp")
	//fmt.Printf("Файл %s успешно скомпилирован.\n", outputFileExe)
	return nil

}

func TestCPPfile(Task_Name string) error {
	path := fmt.Sprintf("src/%v/test.go", Task_Name)
	cmd := exec.Command("go", "run", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	cmd.Run()
	log.Println(string(output))
	return nil
}
