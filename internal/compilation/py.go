package compilation

import (
	"fmt"
	"io/ioutil"
)

func compilePY() {
	// Названия файлов
	baseFile := "1-sum/base.py"
	//sumFile := "sum_py.txt"
	outputFile := "main.py"

	// Чтение содержимого base.py
	baseContent, err := ioutil.ReadFile(baseFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", baseFile, err)
		return
	}

	// Чтение содержимого sum_py.txt
	// sumContent, err := ioutil.ReadFile(sumFile)
	// if err != nil {
	// 	fmt.Printf("Ошибка чтения файла %s: %v\n", sumFile, err)
	// 	return
	// }

	// Создание файла main.py
	err = ioutil.WriteFile(outputFile, append(baseContent), 0644)
	if err != nil {
		fmt.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
		return
	}

	fmt.Printf("Файл %s успешно создан.\n", outputFile)
}
