package filter

type SearchResult struct {
	Index     int
	FileName  string
	Highlight int
	Text      string
}

type SearchOption struct {
	Mode SearchMode
}

type SearchMode int

const (
	Regex SearchMode = iota
	FirstMatch
	WordMatch
	WordMatchIgnoreCase
	FuzzyFind
)

func isValidQuery(q string) bool {
	return false
}

func Search(q string, option *SearchOption) ([]SearchResult, error) {
	return []SearchResult{}, nil
}
