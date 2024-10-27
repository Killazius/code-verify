package compilation

import (
	"compile-server/config"
	"compile-server/internal/models"
	"fmt"
	"os/exec"
)

func MakeFile(path string, lang string, userName string, taskName string) (string, error) {
	if !isValidLang(lang) {
		return "", fmt.Errorf("unsupported language")
	}
	err := config.LoadEnv()
	if err != nil {
		return "", fmt.Errorf("error loading .env file")
	}
	container := config.GetContainer()
	endpoint := config.GetEndpoint()
	containerPath := fmt.Sprintf("%s/%s", container, path)

	userFile := fmt.Sprintf("%s-%s.%s", taskName, userName, lang)

	command := fmt.Sprintf("aws s3 cp --endpoint-url=https://%s s3://%s %s", endpoint, containerPath, userFile)
	cmd := exec.Command(command)
	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf("the file wasn't downloaded from S3 storage")
	}

	return userFile, nil
}

func isValidLang(lang string) bool {
	switch lang {
	case models.LangCpp, models.LangPy:
		return true
	default:
		return false
	}
}
