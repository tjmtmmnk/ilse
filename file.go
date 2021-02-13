package ilse

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	path := path.Join(cfg.userWorkDir, fileName)
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
		log.Fatal("can't open file")
	}
}
