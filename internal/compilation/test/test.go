package test

import (
	"compile-server/internal/compilation"
	"compile-server/internal/handlers/ws/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	Input  [2]string `json:"input"`
	Answer string    `json:"answer"`
}

type testCases []testCase

func Run(command string, args ...string) (string, error) {
	path := "src/1-sum/"
	var tests testCases
	err := readTestCases(path, &tests)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i, test := range tests {
		cmd := exec.CommandContext(ctx, command, args...)
		cmd.Stdin = strings.NewReader(test.Input[0] + "\n" + test.Input[1] + "\n")
		output, errCmd := cmd.CombinedOutput()
		select {
		case <-ctx.Done():
			return utils.Timeout, ctx.Err()
		default:
			if errCmd != nil {
				return "", fmt.Errorf("command execution failed: %w", errCmd)
			}
			result := strings.TrimSpace(string(output))
			if result != test.Answer {
				return fmt.Sprintf("Test case #%d failed. Expected: %s, Got: %s", i+1, test.Answer, result), nil

			}
		}
	}
	return "success", nil
}

func readTestCases(path string, tests *testCases) error {
	file, err := os.Open(fmt.Sprintf("%v%v", path, compilation.TestFile))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &tests)
	if err != nil {
		return err
	}
	return nil
}
