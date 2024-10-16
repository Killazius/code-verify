package compilation

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func CompileCPP(taskName string, userFile string) {
	baseFile := fmt.Sprintf("src/%v/base.cpp", taskName)
	//sumFile := "sum_cpp.txt"
	outputFile := "main.cpp"
	outputFileExe := "main.exe"

	// Чтение содержимого base.cpp
	baseContent, err := ioutil.ReadFile(baseFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return
	}

	// Чтение содержимого sum.txt
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
	fmt.Printf("Файл %s успешно создан.\n", outputFile)

	cmd := exec.Command("g++", "-o", outputFileExe, "main.cpp")

	var out bytes.Buffer
	cmd.Stdout = &out

	err_2 := cmd.Run()
	if err_2 != nil {
		log.Fatal(err)
	}
	fmt.Printf("Файл %s успешно скомпилирован.\n", outputFileExe)
}
