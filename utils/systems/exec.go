package systems

import (
	"errors"
	"os/exec"
)

/*
Exec 执行指定的系统命令并返回其输出结果

参数:

	name string - 要执行的命令名称或路径，不能为空
	args ...string - 传递给命令的参数列表，可变参数

返回值:

	res string - 命令执行的标准输出和错误输出的合并结果
	err error - 执行过程中遇到的错误，包括：
	            - 命令名为空时返回错误
	            - 命令执行失败时返回原始错误

函数流程:
 1. 检查命令名称有效性
 2. 构造命令对象并执行
 3. 捕获输出并处理错误
*/
func Exec(name string, args ...string) (res string, err error) {
	// 检查命令名有效性
	if len(name) == 0 {
		err = errors.New("command name is empty")
		return
	}

	// 创建命令对象并执行，捕获合并输出
	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	res = string(bytes)
	return
}
