package test

import (
	"compile-server/internal/handlers/ws/utils"
	"context"
	"fmt"
	"io"
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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("task_id is incorrect")
	}
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
	if tests == nil {
		return "", fmt.Errorf("no test cases")
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", fmt.Errorf("failed to create stdout pipe: %w", err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return "", fmt.Errorf("failed to create stderr pipe: %w", err)
		}

		if err = cmd.Start(); err != nil {
			return "", fmt.Errorf("failed to start command: %w", err)
		}

		stdoutData, errStdout := io.ReadAll(stdout)
		if errStdout != nil {
			return "", fmt.Errorf("failed to read stdout: %w", errStdout)
		}
		stderrData, errStderr := io.ReadAll(stderr)
		if errStderr != nil {
			return "", fmt.Errorf("failed to read stderr: %w", errStderr)
		}
		if err = cmd.Wait(); err != nil {
			select {
			case <-ctx.Done():
				return "", fmt.Errorf("%v: %w", utils.Timeout, ctx.Err())
			default:
				return "", fmt.Errorf("command execution failed: %w, stderr: %s", err, stderrData)
			}
		}
		result := strings.TrimSpace(string(stdoutData))
		expectedResult := strings.TrimSpace(string(expectedAnswer))
		if result != expectedResult {
			return fmt.Sprintf("Test case #%d failed. Expected: %s, Got: %s", i+1, expectedResult, result), nil
		}
	}
	return utils.OK, nil
}
