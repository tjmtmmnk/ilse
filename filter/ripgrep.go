package filter

import (
	"bufio"
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tjmtmmnk/ilse/util"
)

const (
	timeout = 150 * time.Millisecond
)

type rg struct{}

func newRg() *rg {
	return &rg{}
}

func isValidQuery(q string) bool {
	return q != ""
}

func isValidRegex(q string) bool {
	return true
}

func convert(str string, option *SearchOption) (*SearchResult, error) {
	// first remove reset flag included in path, line
	str = strings.Replace(str, "\x1b[0m", "", 4)
	splitted := strings.Split(str, ":")
	if len(splitted) < 3 {
		return nil, nil
	}

	fileName := splitted[0]
	lineNum, err := strconv.Atoi(splitted[1])
	if err != nil {
		return nil, errors.New("line number wrong format")
	}
	// change reset flag included in text to black foreground
	text := strings.ReplaceAll(splitted[2], "\x1b[0m", "\x1b[39;40m")
	return &SearchResult{fileName, lineNum, text}, nil
}

func (r *rg) Search(q string, option *SearchOption) ([]SearchResult, error) {
	if !isValidQuery(q) {
		return []SearchResult{}, nil
	}

	cmd := []string{
		"rg", "--color=always", "--line-number", "--with-filename",
		"--colors=path:none", "--colors=line:none",
	}

	switch option.Mode {
	case Regex:
		if isValidRegex(q) {
			cmd = append(cmd, "-e")
		} else {
			return []SearchResult{}, nil
		}
	case HeadMatch:
	case WordMatch:
		cmd = append(cmd, "-w")
	}
	if option.Mode != Regex && !option.Case {
		cmd = append(cmd, "-i")
	}

	cmd = append(cmd, q)

	if option.TargetDir != "" {
		cmd = append(cmd, option.TargetDir)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ecmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	stdout, err := ecmd.StdoutPipe()

	if err != nil {
		util.Logger.Warn("command exec stdout pipe error : ", err)
		return []SearchResult{}, err
	}

	if err := ecmd.Start(); err != nil {
		util.Logger.Warn("command exec start error : ", err)
		return []SearchResult{}, err
	}

	results := make([]SearchResult, 0, option.Limit)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if len(results) > option.Limit {
			break
		}
		result, err := convert(scanner.Text(), option)
		if err != nil {
			continue
		}
		if result != nil {
			results = append(results, *result)
		}
	}
	if err := ecmd.Wait(); err != nil {
		util.Logger.Warn("command exec wait error : ", err)
	}

	if ctx.Err() == context.DeadlineExceeded {
		util.Logger.Warn("Timeout")
	}

	return results, nil
}
