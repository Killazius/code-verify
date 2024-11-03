package compilation

import (
	"compile-server/config"
	"compile-server/internal/models"
	"fmt"
	"os"
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
	endpoint := fmt.Sprintf("--endpoint-url=https://%v", config.GetEndpoint())
	container := fmt.Sprintf("s3://%v/%v", config.GetContainer(), path)

	userFile := fmt.Sprintf("%v-%v.%v", taskName, userName, lang)

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

func CreateFile(filePath string, code string, lang string) error {
	if !isValidLang(lang) {
		return fmt.Errorf("unsupported language")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.WriteString(code)

	if err != nil {
		return err
	}
	return nil
}
