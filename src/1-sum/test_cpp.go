package main

import (
	"bufio"
	"bytes"
	"compile-server/internal/models"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	path := "src/1-sum/"
	testFile := models.TestsTxt
	userFileExe := string(os.Args[1])
	userFileExe = fmt.Sprintf("%v%v", path, userFileExe)

	file, err := os.Open(fmt.Sprintf("%v%v", path, testFile))
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)

	var var1, var2 string
	var var3 int

	for scanner.Scan() {
		i := 1
		line := scanner.Text()

		nums := strings.Fields(line)
		if len(nums) == 3 {
			var1 = nums[0]
			var2 = nums[1]
			var3, _ = strconv.Atoi(nums[2])
			cmd := exec.Command("./"+userFileExe, var1, var2)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
				return
			}
			result := out.String()
			cleanedResult := strings.ReplaceAll(result, "\r", "")
			cleanedResult = strings.TrimSpace(cleanedResult)
			resultInt, err := strconv.Atoi(cleanedResult)
			if err != nil {
				fmt.Printf("Ошибка при преобразовании строки '%s' в int: %v\n", cleanedResult, err)
				return
			}
			if resultInt != var3 {
				fmt.Printf("#%v test. Wrong answer", i)
				return
			}
			i++
		} else {
			fmt.Println("Неверный формат строки")
		}
	}
	fmt.Println("Problem solved")
	return
}
