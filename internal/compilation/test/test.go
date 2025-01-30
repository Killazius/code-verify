package test

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

func Run(command string, args ...string) (string, error) {
	path := "src/1-sum/"
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
		cmd := exec.Command(command, args...)
		cmd.Stdin = strings.NewReader(test.Input[0] + "\n" + test.Input[1] + "\n")

		output, errCmd := cmd.CombinedOutput()
		result := strings.TrimSpace(string(output))
		if errCmd != nil {
			return "", errCmd
		}
		if result != test.Answer {
			return fmt.Sprintf("Test case #%d failed.\n", i+1), nil
		}
	}
	return "success", nil
}
