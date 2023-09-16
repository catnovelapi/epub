package newepub

import (
	"fmt"
	"io"
	"log"
	"os"
)

type File struct {
	path string
	file *os.File
}

func Open(path string, mode string) (*File, error) {
	var flag int
	switch mode {
	case "r":
		flag = os.O_RDONLY
	case "w":
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "a":
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil, err
	}

	return &File{path: path, file: f}, nil
}
func NewOpen(path string, mode string) *File {
	f, err := Open(path, mode)
	if err != nil {
		log.Println("Open file error:", err)
		return nil
	}
	return f
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err) // true if file exists
}
func MkdirAll(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Println("创建目录失败:", err)
	}
}
func (f *File) Write(content any) {
	defer f.Close()
	var err error
	switch t := content.(type) {
	case string:
		_, err = f.file.WriteString(t)
	case []byte:
		_, err = f.file.Write(t)
	default:
		err = fmt.Errorf("invalid type: %T", t)
	}
	if err != nil {
		log.Println("file write failed. err: " + err.Error())
	}
}

func (f *File) Read() (string, error) {
	defer f.Close()
	bytes, err := io.ReadAll(f.file)
	return string(bytes), err
}

func (f *File) Close() {
	if err := f.file.Close(); err != nil {
		log.Println("file close failed. err: " + err.Error())
	}
}
