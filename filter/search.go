package filter

import "errors"

type SearchResult struct {
	FileName string
	LineNum  int
	Text     string
}

type SearchOption struct {
	Command   SearchCommand
	Mode      SearchMode
	Case      bool
	TargetDir string
	Limit     int
}

type SearchMode int
type SearchCommand int

type Filter interface {
	Search(string, *SearchOption) ([]SearchResult, error)
}

const (
	Regex SearchMode = iota
	HeadMatch
	WordMatch
	NoneMode
)

const (
	RipGrep SearchCommand = iota
	FuzzySearch
	NoneCommand
)

func NewFilter(option *SearchOption) Filter {
	switch option.Command {
	case RipGrep:
		return newRg()
	case FuzzySearch:
		return newFuzzySearch(option)
	default:
		return newRg()
	}
}

func CommandByName(v string) (SearchCommand, error) {
	switch v {
	case "rg":
		return RipGrep, nil
	case "fuzzy":
		return FuzzySearch, nil
	default:
		return NoneCommand, errors.New("seach command name unmatch")
	}
}

func ModeByName(v string) (SearchMode, error) {
	switch v {
	case "head":
		return HeadMatch, nil
	case "word":
		return WordMatch, nil
	case "regex":
		return Regex, nil
	default:
		return NoneMode, errors.New("search mode name unmatch")
	}
}
