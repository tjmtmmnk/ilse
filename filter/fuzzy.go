package filter

import (
	"bufio"
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
	"github.com/tjmtmmnk/ilse/util"
)

func newFuzzySearch(option *SearchOption) *fuzzySearch {
	fuzzy := &fuzzySearch{}

	gitIgnorePath := filepath.Join(option.TargetDir, ".gitignore")

	if _, err := os.Stat(gitIgnorePath); !os.IsNotExist(err) {
		gitIgnore, err := gitignore.NewGitIgnore(gitIgnorePath)
		if err == nil {
			fuzzy.lookGitIgnore = true
			fuzzy.gitIgnore = gitIgnore
		}
	}
	return fuzzy
}

type fuzzySearch struct {
	cachedFile    map[string]string
	gitIgnore     gitignore.IgnoreMatcher
	lookGitIgnore bool
}

type file struct {
	name string
	text string
}

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
		dir string
	)

	f.cachedFile = make(map[string]string, option.Limit)
	texts := make([]string, 0, option.Limit)
	results := make([]SearchResult, 0, option.Limit)
	isDuplicateLine := make(map[string]bool, option.Limit)

	if option.TargetDir != "" {
		dir = option.TargetDir
	} else {
		dir = "."
	}

	files, err := f.getFiles(dir, option.Limit)
	if err != nil {
		util.Logger.Warn("get file error : ", err)
		return nil, err
	}

	for _, file := range files {
		texts = append(texts, file.text)
	}

	matches := fuzzy.Find(q, texts)
	for _, m := range matches {
		for _, idx := range f.reIndex(m.MatchedIndexes, len(q)) {
			if len(files) <= m.Index {
				continue
			}
			fileName := files[m.Index].name
			lineNum, text := f.getLine(files[m.Index].name, idx+len(q))
			if text == "" {
				continue
			}
			key := fmt.Sprintf("%s%d", fileName, lineNum)
			if isDuplicateLine[key] {
				continue
			}
			isDuplicateLine[key] = true
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
	to := pos
	if to > len(text) {
		to = len(text) - 1
	}
	lineNum := 1

	var sb strings.Builder
	sb.Grow(to)
	for _, c := range text[:to] {
		sb.WriteRune(c)

		if c == '\n' {
			lineNum++
			sb.Reset()
		}
	}
	return lineNum, sb.String()
}

func (f *fuzzySearch) getFiles(dir string, limit int) ([]file, error) {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	files := make([]file, 0, limit*2)

	fileNames, err := f.getFileNames(dir, limit)
	if err != nil {
		util.Logger.Warn("get file name error : ", err)
		return nil, err
	}

	for _, name := range fileNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var sb strings.Builder
			sb.Grow(100)
			fp, err := os.Open(name)
			if err != nil {
				util.Logger.Warn("file open error : ", err)
				return
			}
			defer fp.Close()

			scanner := bufio.NewScanner(fp)

			for scanner.Scan() {
				if scanner.Err() != nil {
					util.Logger.Warn("scan error : ", err)
					return
				}
				if scanner.Text() != "" {
					sb.WriteString(scanner.Text())
				}
			}
			text := sb.String()
			mimeType := http.DetectContentType([]byte(text))
			isBinary := !strings.HasPrefix(mimeType, "text/plain")
			if isBinary {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			files = append(files, file{name, text})
		}(name)
	}
	wg.Wait()
	return files, nil
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
	fileNames := make([]string, 0, limit*2)
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
