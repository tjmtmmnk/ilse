package util

import (
	"fmt"
	"os"
	"os/exec"
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
	isManaged, err := isManagedByGit()
	if err != nil {
		return ".", nil
	}
	if isManaged {
		repo, err := getGitRepository()
		if err != nil {
			return ".", err
		}
		if repo != "" {
			return repo, nil
		}
	}
	return ".", nil
}

func getGitRepository() (string, error) {
	cmd := []string{"git", "rev-parse", "--show-toplevel"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	repo := string(out)
	return strings.TrimRight(repo, "\r\n"), nil
}

func isManagedByGit() (bool, error) {
	cmd := []string{"git", "rev-parse", "--all"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return false, err
	}
	isManaged := len(out) != 0
	return isManaged, nil
}
