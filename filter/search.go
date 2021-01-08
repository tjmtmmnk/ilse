package filter

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
)

const (
	RipGrep SearchCommand = iota
	FuzzySearch
)

func NewFilter(cmd SearchCommand) Filter {
	switch cmd {
	case RipGrep:
		return newRg()
	case FuzzySearch:
		return newFuzzySearch()
	default:
		return newRg()
	}
}

func DefaultOption() *SearchOption {
	return &SearchOption{
		Command: RipGrep,
		Mode:    HeadMatch,
		Case:    false,
	}
}
