package compilation

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	fmt.Printf("Файл %s успешно создан.\n", outputFile)

	err = CompileCPPfile(outputFileExe, outputFile)
	if err != nil {
		log.Println("Файл не скомпилирован.")
		return
	}

	err = TestCPPfile(outputFileExe, taskName)
	if err != nil {
		log.Println("Тесты не прошли")
		return
	}

}

func CompileCPPfile(outputFileExe string, outputFile string) error {
	cmd := exec.Command("g++", "-o", outputFileExe, outputFile)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	os.Remove("user.cpp")
	os.Remove("solution.cpp")
	fmt.Printf("Файл %s успешно скомпилирован.\n", outputFileExe)
	return nil

}

func TestCPPfile(outputFileExe string, Task_Name string) error {
	testsPath := fmt.Sprintf("src/%v/test.txt", Task_Name)
	file, err := os.Open(testsPath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var var1, var2 string
	var var3 int

	for scanner.Scan() {
		line := scanner.Text()

		nums := strings.Fields(line)
		if len(nums) == 3 {
			var1 = nums[0]
			var2 = nums[1]
			var3, _ = strconv.Atoi(nums[2])
			cmd := exec.Command("./"+outputFileExe, var1, var2)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
				return err
			}
			result := out.String()
			cleanedResult := strings.ReplaceAll(result, "\r", "")
			cleanedResult = strings.TrimSpace(cleanedResult)
			resultInt, err := strconv.Atoi(cleanedResult)
			if err != nil {
				fmt.Printf("Ошибка при преобразовании строки '%s' в int: %v\n", cleanedResult, err)
				return err
			}
			fmt.Printf("Вызывается тест. Данные: %d - %d\n", resultInt, var3)
			if resultInt != var3 {
				fmt.Println("Тест провален!")
				return nil
			}

		} else {
			fmt.Println("Неверный формат строки")
		}
	}
	fmt.Println("Задача решена верно!")
	os.Remove(outputFileExe)
	return nil
}
