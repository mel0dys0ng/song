package files

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Exist return the file or dir exist
func Exist(path string) (res bool) {
	res = true
	if _, e := os.Stat(path); e != nil {
		res = false
	}
	return
}

// Abs 返回绝对路径。若发生错误，则返回空。
func Abs(path string) string {
	abspath, _ := filepath.Abs(path)
	return abspath
}

// ReadFile 读文件内容
func ReadFile(path string) (content string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	for {
		data, err := reader.ReadString('\n')
		content += data
		if err == io.EOF {
			break
		}
	}

	return
}

// ReadYamlFile 读取Yaml文件内容并解析至data
func ReadYamlFile(path string, data any) (err error) {
	content, err := os.ReadFile(path)
	if err == nil && len(content) > 0 {
		err = yaml.Unmarshal(content, data)
	}
	return
}

// WriteFile 写文件内容
func WriteFile(path, content string, append bool, mode fs.FileMode) (
	n int, err error) {

	file, err := FileWriter(path, append, mode)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(file)
	for i := 0; i < 3; i++ {
		n, err = writer.WriteString(content)
		if err == nil {
			break
		}
	}

	err = writer.Flush()

	return
}

// FileWriter 返回文件句柄，若文件目录不存在，则自动创建
func FileWriter(path string, append bool, mode fs.FileMode) (
	file *os.File, err error) {
	fp, err := MakeDir(path, mode)
	if err != nil {
		return
	}

	appendOrTrunc := os.O_TRUNC
	if append {
		appendOrTrunc = os.O_APPEND
	}

	return os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|appendOrTrunc, mode)
}

// MakeDir make dir all
func MakeDir(path string, mode fs.FileMode) (fp string, err error) {
	fp, err = filepath.Abs(path)
	if err != nil {
		return
	}

	if _, e := os.Stat(fp); e != nil {
		dir, _ := filepath.Split(fp)
		if _, er := os.Stat(dir); er != nil {
			if err = os.MkdirAll(dir, mode); err == nil {
				err = os.Chmod(dir, mode)
			}
			if err != nil {
				return
			}
		}
	}

	return
}

func ReadDir(path string) (res []fs.FileInfo, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	res = make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, er := entry.Info()
		if err = er; err != nil {
			return
		}
		res = append(res, info)
	}

	return
}

// WalkDir 遍历当前目录
func WalkDir(path string, walk func(path string, info fs.FileInfo) bool) (
	result []string, err error) {
	//获取文件或目录相关信息
	list, err := ReadDir(path)
	if err != nil {
		return
	}

	for i := range list {
		if walk == nil || walk(list[i].Name(), list[i]) {
			result = append(result, path+"/"+list[i].Name())
		}
	}

	return
}

// WalkDirs 迭代遍历当前目录以及子目录
func WalkDirs(path string, walk func(path string, info fs.FileInfo) bool) (
	result []string, err error) {
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if walk == nil || walk(path, info) {
			result = append(result, path)
		}
		return err
	})
	return
}
