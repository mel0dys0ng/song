package systems

import (
	"errors"
	"os/exec"
)

// Exec 执行命令，name：命令名称，args：名称参数
func Exec(name string, args ...string) (res string, err error) {
	if len(name) == 0 {
		err = errors.New("command name is empty")
		return
	}

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	res = string(bytes)
	return
}
