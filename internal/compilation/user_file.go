package compilation

import (
	"fmt"
	"os/exec"
)

func MakeFile(Link string, format string) error {

	endpoint := "--endpoint-url=https://s3.ru-1.storage.selcloud.ru"
	cmd := exec.Command("aws", "s3", "cp")
	cmd.Args = append(cmd.Args, endpoint, Link, fmt.Sprintf("user.%s", format))
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("ошибка выполнения команды: %v", err)
	}

	return nil
}
