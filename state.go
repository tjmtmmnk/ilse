package ilse

import (
	"sync"

	"github.com/tjmtmmnk/ilse/filter"
)

type state struct {
	mutex        sync.RWMutex
	currentPage  string
	matched      []filter.SearchResult
	fileCache    map[string][]string
	targetDir    []string
	ignoreDir    []string
	searchOption *filter.SearchOption
}

func newState() *state {
	return &state{
		searchOption: &filter.SearchOption{
			Command: filter.RipGrep,
			Mode:    filter.FirstMatch,
		},
	}
}
