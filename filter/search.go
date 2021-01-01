package filter

type SearchResult struct {
	FileName string
	LineNum  int
	Text     string
}

type SearchOption struct {
	Command SearchCommand
	Mode    SearchMode
}

type SearchMode int
type SearchCommand int

type filter interface {
	Search(string, *SearchOption) ([]SearchResult, error)
}

const (
	Regex SearchMode = iota
	FirstMatch
	FirstMatchCase
	WordMatch
	WordMatchCase
)

const (
	RipGrep SearchCommand = iota
	FuzzySearch
)

func NewFilter(cmd SearchCommand) filter {
	switch cmd {
	case RipGrep:
		return &rg{}
	case FuzzySearch:
		return &fuzzySearch{}
	default:
		return &rg{}
	}
}
