package filter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
)

func newFuzzySearch() *fuzzySearch {
	return &fuzzySearch{
		cachedFile:  make(map[string]string),
		isDuplicate: make(map[string]bool),
	}
}

type fuzzySearch struct {
	cachedFile  map[string]string
	isDuplicate map[string]bool
}

type file struct {
	name string
	text string
}

func (f *fuzzySearch) purge() {
	f.cachedFile = make(map[string]string)
	f.isDuplicate = make(map[string]bool)
}

func (f *fuzzySearch) Search(q string, option *SearchOption) ([]SearchResult, error) {
	var (
		results []SearchResult
		texts   []string
	)
	f.purge()

	length := len(q)

	files := f.getFiles()

	for _, file := range files {
		texts = append(texts, file.text)
	}

	matches := fuzzy.Find(q, texts)
	for _, m := range matches {
		for _, idx := range f.reIndex(m.MatchedIndexes, length) {
			fileName := files[m.Index].name
			lineNum, text := f.getLine(files[m.Index].name, idx+length)
			if text == "" {
				continue
			}
			key := fmt.Sprintf("%s%d", fileName, lineNum)
			if f.isDuplicate[key] {
				continue
			}
			f.isDuplicate[key] = true
			results = append(results, SearchResult{fileName, lineNum, text})
		}
	}

	return results, nil
}

func (f *fuzzySearch) getLine(fileName string, pos int) (int, string) {
	var text string
	if f.cachedFile[fileName] != "" {
		text = f.cachedFile[fileName]
	} else {
		content, _ := ioutil.ReadFile(fileName)
		text = string(content)
		f.cachedFile[fileName] = text
	}
	lineNum := 1
	lineText := ""
	to := pos
	if to > len(text) {
		to = len(text) - 1
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

func (f *fuzzySearch) getFiles() []file {
	var wg sync.WaitGroup
	var mu sync.Mutex
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
			mimeType := http.DetectContentType(text)
			isBinary := !strings.HasPrefix(mimeType, "text/plain")
			if isBinary {
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
func (f *fuzzySearch) reIndex(indexes []int, queryLength int) []int {
	if len(indexes) == 0 {
		return []int{}
	}
	baseIdx := indexes[0]
	reIndexes := []int{baseIdx}
	for _, idx := range indexes {
		isSameWord := idx <= baseIdx+queryLength
		if isSameWord {
			continue
		}
		reIndexes = append(reIndexes, idx)
		baseIdx = idx
	}
	return reIndexes
}

func (f *fuzzySearch) getFileNames() ([]string, error) {
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
