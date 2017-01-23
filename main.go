package main

import (
	"os"

	"github.com/ohtomi/scrapbox/command"
)

func main() {
	command.InitializeMeta()
	os.Exit(Run(os.Args[1:]))
}
