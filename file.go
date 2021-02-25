package ilse

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tjmtmmnk/ilse/util"
)

func getEditorLineFlag() string {
	editor := os.Getenv("EDITORLINEFLAG")
	if len(editor) == 0 {
		editor = "+"
	}
	return editor
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}
	return editor
}

func openFile(fileName string, lineNum int) {
	editor := getEditor()

	path := filepath.Join(conf.userWorkDir, fileName)
	if filepath.IsAbs(fileName) {
		path = fileName
	}
	lineFlag := fmt.Sprintf("%s%d", getEditorLineFlag(), lineNum)

	var cmd *exec.Cmd
	_, file := filepath.Split(editor)
	switch file {
	case "emacs", "emacsclient":
		// emacs expects the line before the file
		cmd = exec.Command(editor, lineFlag, path)
	default:
		// default to adding the line after the file
		cmd = exec.Command(editor, path, lineFlag)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		util.Logger.Error("can't open file")
	}
}
