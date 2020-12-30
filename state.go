package ilse

import "github.com/tjmtmmnk/ilse/filter"

type state struct {
	currentPage  string
	matched      []filter.SearchResult
	fileCache    map[string][]string
	targetDir    []string
	ignoreDir    []string
	searchOption filter.SearchOption
}
