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
	endpoint := fmt.Sprintf("--endpoint-url=https://%s", config.GetEndpoint())
	container := fmt.Sprintf("s3://%s/%s", config.GetContainer(), path)

	userFile := fmt.Sprintf("%s-%s.%s", taskName, userName, lang)

	cmd := exec.Command("aws", "s3", "cp")
	cmd.Args = append(cmd.Args, endpoint, container, userFile)
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
