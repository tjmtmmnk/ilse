package main

import (
	"log"

	"github.com/integrii/flaggy"
	"github.com/tjmtmmnk/ilse"
)

func initFlags(config *ilse.Config) {
	const (
		name        = "ilse"
		description = "ilse is TUI grep tool like IntelliJ"
		version     = "0.1"
	)
	flaggy.SetVersion(version)
	flaggy.SetName(name)
	flaggy.SetDescription(description)

	flaggy.Int(&config.MaxSearchResults, "m", "max-search-results", "Max number of search results")
	flaggy.String(&config.SearchCommand, "f", "filter", "select filter ('rg', 'fuzzy')")
	flaggy.String(&config.SearchMode, "fm", "filter-mode", "select filter mode ('head', 'word', 'regex')")
	flaggy.String(&config.Theme, "t", "preview-theme", "select bat theme for preview")
	flaggy.Bool(&config.CaseSensitive, "c", "case", "case sensitive")

	flaggy.Parse()

}

func main() {
	config, err := ilse.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	initFlags(config)
	if err := ilse.Init(config); err != nil {
		log.Fatal(err)
	}
	if err := ilse.Run(); err != nil {
		log.Fatal(err)
	}
}
