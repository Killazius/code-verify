package test

import (
	"compile-server/internal/handlers/ws/utils"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	InputFile  string
	AnswerFile string
}

type testCases []testCase

func readTestCases(path string) (testCases, error) {
	var tests testCases

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	inputFiles := make(map[string]string)
	answerFiles := make(map[string]string)

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "in") {
			num := strings.TrimPrefix(name, "in")
			inputFiles[num] = filepath.Join(path, name)
		} else if strings.HasPrefix(name, "out") {
			num := strings.TrimPrefix(name, "out")
			answerFiles[num] = filepath.Join(path, name)
		}
	}

	for num, inputFile := range inputFiles {
		answerFile, exists := answerFiles[num]
		if !exists {
			return nil, fmt.Errorf("missing answer file for input file: %v", inputFile)
		}
		tests = append(tests, testCase{InputFile: inputFile, AnswerFile: answerFile})
	}

	return tests, nil
}

func Run(command, taskName string, args ...string) (string, error) {
	path := fmt.Sprintf("src/%v", taskName)
	tests, errTests := readTestCases(path)
	if errTests != nil {
		return "", errTests
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i, test := range tests {
		inputData, err := os.ReadFile(test.InputFile)
		if err != nil {
			return "", fmt.Errorf("failed to read input file: %w", err)
		}
		expectedAnswer, err := os.ReadFile(test.AnswerFile)
		if err != nil {
			return "", fmt.Errorf("failed to read answer file: %w", err)
		}
		cmd := exec.CommandContext(ctx, command, args...)
		cmd.Stdin = strings.NewReader(string(inputData))
		output, errCmd := cmd.CombinedOutput()
		select {
		case <-ctx.Done():
			return utils.Timeout, ctx.Err()
		default:
			if errCmd != nil {
				return "", fmt.Errorf("command execution failed: %w", errCmd)
			}
			result := strings.TrimSpace(string(output))
			expectedResult := strings.TrimSpace(string(expectedAnswer))
			if result != expectedResult {
				return fmt.Sprintf("Test case #%d failed. Expected: %s, Got: %s", i+1, expectedResult, result), nil
			}
		}
	}
	return utils.OK, nil
}
