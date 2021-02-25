package ilse

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/util"
)

var (
	tree *tview.TreeView
)

func initTree() error {
	rootDir, err := util.GetUserWorkDir()
	if err != nil {
		return err
	}

	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)

	addNode(root, rootDir)

	tree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentPage, _ := pages.GetFrontPage()
		if currentPage != treePage {
			return event
		}
		if event.Key() == tcell.KeyRight || event.Rune() == 'l' {
			if err := expand(tree.GetCurrentNode()); err != nil {
				util.Logger.Error("expand error : %v", err)
			}
		}
		return event
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		path := rootDir
		if reference != nil {
			path = reference.(string)
		}
		app.searchOption.TargetDir = path
		pages.SwitchToPage(mainPage)
	})

	tree.SetDoneFunc(func(key tcell.Key) {
		pages.SwitchToPage(mainPage)
	})

	return nil
}

func addNode(target *tview.TreeNode, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(file.IsDir())
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
	return nil
}

func expand(node *tview.TreeNode) error {
	reference := node.GetReference()
	if reference == nil {
		return nil
	}
	children := node.GetChildren()
	if len(children) == 0 {
		path := reference.(string)
		if err := addNode(node, path); err != nil {
			return err
		}
	} else {
		node.SetExpanded(!node.IsExpanded())
	}
	return nil
}
