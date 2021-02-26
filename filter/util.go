package filter

import (
	"errors"
	"strconv"
	"strings"
)

func isValidQuery(q string) bool {
	return q != ""
}

func isValidRegex(q string) bool {
	return true
}

func convert(result string, option *SearchOption) ([]SearchResult, error) {
	results := make([]SearchResult, 0, option.Limit)
	for i, s := range strings.Split(string(result), "\n") {
		if i > option.Limit {
			break
		}
		result, err := split(s, option)
		if err != nil {
			return []SearchResult{}, err
		}
		if result != nil {
			results = append(results, *result)
		}
	}
	return results, nil
}

func split(str string, option *SearchOption) (*SearchResult, error) {
	switch option.Command {
	case RipGrep:
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
	return nil, errors.New("wrong option")
}
