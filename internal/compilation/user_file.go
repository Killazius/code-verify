package compilation

import (
	"fmt"
	"os/exec"
)

func MakeFile(Link string, format string) {
	endpoint := "--endpoint-url=https://s3.ru-1.storage.selcloud.ru"
	url := fmt.Sprintf("aws %v s3 cp  %v user.%v", endpoint, Link, format)
	exec.Command(url)
}
