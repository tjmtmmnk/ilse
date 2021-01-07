package ilse

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/tjmtmmnk/ilse/filter"
)

type ilse struct {
	state        *state
	config       *config
	searchOption *filter.SearchOption
	screen       tcell.Screen
}

var (
	app *ilse
)

func initApp() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	state := newState()
	config, err := newConfig()
	if err != nil {
		return err
	}
	searchOption := filter.DefaultOption()
	app = &ilse{
		screen:       screen,
		state:        state,
		config:       config,
		searchOption: searchOption,
	}
	return nil
}

func Init() error {
	if err := initApp(); err != nil {
		return err
	}
	initFrame()
	return nil
}

func Run() error {
	return frame.Run()
}
