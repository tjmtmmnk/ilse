package ilse

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/tjmtmmnk/ilse/filter"
)

type ilse struct {
	state        *state
	searchOption *filter.SearchOption
	screen       tcell.Screen
}

var (
	app    *ilse
	cfg    *config
	logger *Logger
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

	os.MkdirAll(cfg.homeDir, 0766)
	searchOption := filter.DefaultOption()
	searchOption.Limit = cfg.maxSearchResults
	app = &ilse{
		screen:       screen,
		state:        state,
		searchOption: searchOption,
	}
	return nil
}

func Init() error {
	var err error
	cfg, err = newConfig()
	if err != nil {
		return err
	}

	if err := initApp(); err != nil {
		return err
	}
	if err := initFrame(); err != nil {
		return err
	}

	logger, err = newLogger()
	if err != nil {
		return err
	}

	return nil
}

func Run() error {
	return frame.Run()
}
