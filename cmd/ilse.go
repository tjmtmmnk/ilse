package main

import (
	"log"

	"github.com/tjmtmmnk/ilse"
)

func main() {
	if err := ilse.Init(); err != nil {
		log.Fatal(err)
	}
	if err := ilse.Run(); err != nil {
		log.Fatal(err)
	}
}
