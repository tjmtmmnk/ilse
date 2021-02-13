package filter

import (
	"strconv"
	"strings"
)

func isValidQuery(q string) bool {
	return q != ""
}

func isValidRegex(q string) bool {
	return true
}

func convert(result string, option *SearchOption) []SearchResult {
	var results []SearchResult
	for i, s := range strings.Split(string(result), "\n") {
		if i > option.Limit {
			break
		}
		ignore, result := split(s, option)
		if !ignore {
			results = append(results, result)
		}
	}
	return results
}

func split(str string, option *SearchOption) (ignore bool, result SearchResult) {
	switch option.Command {
	case RipGrep:
		// first remove reset flag included in path, line
		str = strings.Replace(str, "\x1b[0m", "", 4)
		splitted := strings.Split(str, ":")
		if len(splitted) < 3 {
			ignore = true
			return
		}

		fileName := splitted[0]
		lineNum, _ := strconv.Atoi(splitted[1])
		// change reset flag included in text to black foreground
		text := strings.ReplaceAll(splitted[2], "\x1b[0m", "\x1b[39;40m")
		result = SearchResult{fileName, lineNum, text}
	}
	return
}
