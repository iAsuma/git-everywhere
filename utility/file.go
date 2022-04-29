package utility

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetPwd 获取当前路径
func GetPwd() string {
	dir, err := os.Getwd() //当前的目录
	if err != nil {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Println("can not get current path")
		}
	}
	return dir
}

func readFileLines(fileName string, fc func(text string) (err error)) (err error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}

		if err = fc(line); err != nil {
			return err
		}
	}
	return nil
}
