package util

import (
	"fmt"
	"os"
	"os/exec"
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

func GetHomeDir() string {
	userHomeDir, _ := os.UserHomeDir()
	homeDir := fmt.Sprintf("%s/.ilse", userHomeDir)
	return homeDir
}

// if use git, return repository
// else return current directory
func GetUserWorkDir() (string, error) {
	repo, err := getGitRepository()
	if err != nil {
		return "", err
	}
	if repo != "" {
		return repo, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	evaled, err := filepath.EvalSymlinks(wd)
	if err != nil {
		return "", err
	}

	return evaled, nil
}

func getGitRepository() (string, error) {
	cmd := []string{"git", "rev-parse", "--show-toplevel"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	str := strings.Split(string(out), "\n")
	if len(str) == 0 {
		return "", nil
	}
	repo := str[0]
	return repo, nil
}
