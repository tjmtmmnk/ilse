package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ShortFileName(fileName string) string {
	sp := strings.Split(fileName, "/")
	if len(sp) == 1 {
		return sp[0]
	} else if len(sp) == 2 {
		return fmt.Sprintf("%c/%s", sp[0][0], sp[1])
	} else {
		last := len(sp) - 1
		return fmt.Sprintf("%c/%c/%s", sp[last-2][0], sp[last-1][0], sp[last])
	}
}

func GetWorkDir() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	evaled, err := filepath.EvalSymlinks(workDir)
	if err != nil {
		return "", err
	}

	return evaled, nil
}
