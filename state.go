package ilse

import (
	"sync"

	"github.com/tjmtmmnk/ilse/filter"
)

type state struct {
	mutex       sync.RWMutex
	currentPage string
	matched     []filter.SearchResult
	fileCache   map[string][]string
	targetDir   []string
	ignoreDir   []string
}

func newState() *state {
	return &state{}
}
