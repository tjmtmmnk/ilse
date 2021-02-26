package filter

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/monochromegane/go-gitignore"
	"github.com/sahilm/fuzzy"
)

func newFuzzySearch(option *SearchOption) *fuzzySearch {
	fs := &fuzzySearch{
		cachedFile:  make(map[string]string),
		isDuplicate: make(map[string]bool),
	}
	fs := &fuzzySearch{}
	gitIgnorePath := filepath.Join(option.TargetDir, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); !os.IsNotExist(err) {
		gitIgnore, err := gitignore.NewGitIgnore(gitIgnorePath)
		if err == nil {
			fs.lookGitIgnore = true
			fs.gitIgnore = gitIgnore
		}
	}
	return fs
}

type fuzzySearch struct {
	cachedFile    map[string]string
	isDuplicate   map[string]bool
	gitIgnore     gitignore.IgnoreMatcher
	lookGitIgnore bool
}

type file struct {
	name string
	text string
}

func (f *fuzzySearch) purge() {
	f.cachedFile = make(map[string]string)
	f.isDuplicate = make(map[string]bool)
type fileSystem struct {
	fs.ReadDirFS
}

func (fs *fileSystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (fs *fileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (f *fuzzySearch) Search(q string, option *SearchOption) ([]SearchResult, error) {
	var (
		results []SearchResult
		texts   []string
		dir     string
	)
	f.purge()

	length := len(q)

	if option.TargetDir != "" {
		dir = option.TargetDir
	} else {
		dir = "."
	}
	files := f.getFiles(dir)

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

func (f *fuzzySearch) getFiles(dir string) []file {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var files []file

	fileNames, _ := f.getFileNames(dir)
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
			defer mu.Unlock()
			files = append(files, file{name, string(text)})
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

func (f *fuzzySearch) getFileNames(root string, limit int) ([]string, error) {
	fsm := &fileSystem{}
	fileNames := make([]string, 0, limit)
	err := fs.WalkDir(fsm, root, func(path string, d fs.DirEntry, err error) error {
		if f.lookGitIgnore && f.gitIgnore.Match(path, d.IsDir()) {
			return fs.SkipDir
		}

		// skip hidden directory or file
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
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
