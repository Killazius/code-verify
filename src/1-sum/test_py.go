package main

import (
	"compile-server/internal/compilation"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	Input  [2]string `json:"input"`
	Answer string    `json:"answer"`
}

type testCases []testCase

func main() {
	path := "src/1-sum/"
	solutionFile := os.Args[1]

	var tests testCases
	file, err := os.Open(fmt.Sprintf("%v%v", path, compilation.TestFile))
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteValue, &tests)
	if err != nil {
		log.Fatal(err)
	}

	for i, test := range tests {
		num1, num2 := test.Input[0], test.Input[1]
		cmd := exec.Command("python3", solutionFile, num1, num2)
		output, errCmd := cmd.CombinedOutput()
		result := strings.TrimSpace(string(output))
		if errCmd != nil {
			log.Fatal(errCmd)
			return
		}
		if result != test.Answer {
			fmt.Printf("Test case #%d failed.\n", i+1)
			return
		}
	}
	fmt.Println("success")
}
