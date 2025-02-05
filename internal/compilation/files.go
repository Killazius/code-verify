package compilation

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/exec"
)

type Lang string

const (
	LangCpp Lang = "cpp"
	LangPy  Lang = "py"
	LangGO  Lang = "go"
)

// MakeFile Dispatched: убрали функцию сохранения на s3. думать как сделать s3 для хранения попыток и всего прочего
func _(path string, lang Lang, userName string, taskName string) (string, error) {
	if !isValidLang(lang) {
		return "", fmt.Errorf("unsupported language")
	}
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file")
	}
	endpoint := fmt.Sprintf("--endpoint-url=https://%v", "config.GetEndpoint()")
	container := fmt.Sprintf("s3://%v/%v", "config.GetContainer()", path)

	userFile := fmt.Sprintf("%v-%v.%v", taskName, userName, lang)

	cmd := exec.Command("aws", "s3", "cp")
	cmd.Args = append(cmd.Args, endpoint, container, userFile)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("the file wasn't downloaded from S3 storage")
	}

	return userFile, nil
}

func isValidLang(lang Lang) bool {
	switch lang {
	case LangCpp, LangPy, LangGO:
		return true
	default:
		return false
	}
}

func CreateFile(filePath string, code string, lang Lang) error {
	if !isValidLang(lang) {
		return fmt.Errorf("unsupported language")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			return
		}
	}(file)

	_, err = file.WriteString(code)

	if err != nil {
		return err
	}
	return nil
}
