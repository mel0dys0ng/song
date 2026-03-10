package fs

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Exists returns whether the file or directory exists
func Exists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	return
}

// Abs returns the absolute path
func Abs(path string) (abspath string, err error) {
	return filepath.Abs(path)
}

// ReadFile reads file content
func ReadFile(path string) (content string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}

	content = string(data)
	return
}

// ReadYamlFile reads and parses YAML file content
func ReadYamlFile(path string, data any) (err error) {
	content, err := os.ReadFile(path)
	if err == nil && len(content) > 0 {
		err = yaml.Unmarshal(content, data)
	}
	return
}

// WriteFile writes content to file
func WriteFile(path, content string, append bool, mode fs.FileMode) (
	n int, err error) {

	file, err := FileWriter(path, append, mode)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < 3; i++ {
		n, err = writer.WriteString(content)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err == nil {
		err = writer.Flush()
	}

	return
}

// FileWriter returns file handle, automatically creates directory if not exists
func FileWriter(path string, append bool, mode fs.FileMode) (file *os.File, err error) {
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

// MakeDir creates directories recursively
func MakeDir(path string, mode fs.FileMode) (fp string, err error) {
	fp, err = filepath.Abs(path)
	if err != nil {
		return
	}

	dir, _ := filepath.Split(fp)
	if dir != "" {
		if err = os.MkdirAll(dir, mode); err != nil {
			return
		}
	}

	return
}

// ReadDir reads directory contents
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

// WalkDir walks through the current directory
func WalkDir(path string, walk func(name string, info fs.FileInfo) bool) (
	result []string, err error) {
	//获取文件或目录相关信息
	list, err := ReadDir(path)
	if err != nil {
		return
	}

	result = make([]string, 0, len(list))
	for i := range list {
		if walk == nil || walk(list[i].Name(), list[i]) {
			result = append(result, filepath.Join(path, list[i].Name()))
		}
	}

	return
}

// WalkDirs walks through current directory and subdirectories recursively
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
