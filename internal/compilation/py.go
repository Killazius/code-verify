package compilation

import (
	"fmt"
	"log"
	"os"
)

func MakePYfile(taskName string, userFile string) {
	pathTask := fmt.Sprintf("src/%v", taskName)
	baseFile := fmt.Sprintf("%v/base.py", pathTask)
	outputFile := fmt.Sprintf("%v/solution.py", pathTask)

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

	err = os.WriteFile(outputFile, append(baseContent, userContent...), 0644)
	if err != nil {
		log.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
		return
	}

	//err = CompileCPPfile(userContent, outputFile, taskName)
	//if err != nil {
	//	log.Printf("Файл не скомпилирован.")
	//	return
	//}
	//err = TestCPPfile(taskName)
	//if err != nil {
	//	log.Printf("Ошибка во время тестирования: %v", err)
	//	return
	//}

}
