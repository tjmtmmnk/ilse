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
	app  *ilse
	conf *Config
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

	os.MkdirAll(conf.homeDir, 0766)

	command, err := filter.CommandByName(conf.SearchCommand)
	if err != nil {
		return err
	}
	mode, err := filter.ModeByName(conf.SearchMode)
	if err != nil {
		return err
	}
	searchOption := &filter.SearchOption{
		Command:   command,
		Mode:      mode,
		Case:      conf.CaseSensitive,
		Limit:     conf.MaxSearchResults,
		TargetDir: conf.userWorkDir,
	}

	app = &ilse{
		screen:       screen,
		state:        state,
		searchOption: searchOption,
	}
	return nil
}

func Init(cfg *Config) error {
	conf = cfg

	if err := initApp(); err != nil {
		return err
	}
	if err := initFrame(); err != nil {
		return err
	}

	return nil
}

func Run() error {
	return frame.Run()
}
