package compilation

import (
	"fmt"
	"os/exec"
)

const (
	endpoint = "--endpoint-url=https://s3.ru-1.storage.selcloud.ru"
)

func MakeFile(path string, lang string, userName string, taskName string) (string, error) {

	container := "s3://container-studying-2/"
	container += path

	userFile := fmt.Sprintf("%s-%s.%s", taskName, userName, lang)
	cmd := exec.Command("aws", "s3", "cp")
	cmd.Args = append(cmd.Args, endpoint, container, userFile)
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("the file wasn't downloaded: %v", err)
	}

	return userFile, nil
}
