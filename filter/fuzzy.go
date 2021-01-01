package filter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/sahilm/fuzzy"
)

type Fuzzy struct{}

type searchResultWithScore struct {
	*SearchResult
	score int
}

type file struct {
	name string
	text string
}

func (f *Fuzzy) Search(q string, option *SearchOption) ([]SearchResult, error) {
	var (
		results          []SearchResult
		resultsWithScore []searchResultWithScore
		texts            []string
	)

	files := f.getFiles()

	for _, file := range files {
		texts = append(texts, file.text)
	}

	matches := fuzzy.Find(q, texts)
	for _, m := range matches {
		for _, idx := range f.reIndex(m.MatchedIndexes, len(q)) {
			fileName := files[m.Index].name
			lineNum, text := f.getLine(files[m.Index].name, idx+len(q))
			if text == "" {
				continue
			}
			result := SearchResult{fileName, lineNum, text}
			resultsWithScore = append(resultsWithScore, searchResultWithScore{&result, m.Score})
		}
	}
	sort.Slice(resultsWithScore, func(i, j int) bool {
		return resultsWithScore[i].score > resultsWithScore[j].score
	})

	for _, res := range resultsWithScore {
		results = append(results, *res.SearchResult)
	}

	return results, nil
}

func (f *Fuzzy) getLine(fileName string, pos int) (int, string) {
	text, _ := ioutil.ReadFile(fileName)
	lineNum := 1
	lineText := ""
	to := pos
	if pos > len(text) {
		pos = len(text) - 1
	}
	for _, c := range text[:to] {
		lineText += string(c)

		if c == '\n' {
			lineNum++
			lineText = ""
		}
	}
	return lineNum, lineText
}

func (f *Fuzzy) getFiles() []file {
	var wg sync.WaitGroup
	var mu sync.RWMutex
	var files []file
	fileNames, _ := f.getFileNames()
	for _, name := range fileNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			text, err := ioutil.ReadFile(name)
			if err != nil {
				return
			}
			mu.Lock()
			files = append(files, file{name, string(text)})
			mu.Unlock()
		}(name)
	}
	wg.Wait()
	return files
}

func (f *Fuzzy) reIndex(indexes []int, queryLength int) []int {
	if len(indexes) == 0 {
		return []int{}
	}
	baseIdx := indexes[0]
	reIndexes := []int{baseIdx}
	for _, idx := range indexes {
		isSameIndex := idx <= baseIdx+queryLength
		if isSameIndex {
			continue
		}
		reIndexes = append(reIndexes, idx)
		baseIdx = idx
	}
	return reIndexes
}

func (f *Fuzzy) getFileNames() ([]string, error) {
	var fileNames []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileNames = append(fileNames, path)
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return fileNames, nil
}
