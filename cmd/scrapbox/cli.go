package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/ohtomi/scrapbox/cmd/scrapbox/command"
)

func Run(args []string) int {

	// Meta-option for executables.
	// It defines output color and its stdout/stderr stream.
	meta := &command.Meta{
		Ui: &cli.ColoredUi{
			OutputColor: cli.UiColorNone,
			InfoColor:   cli.UiColorBlue,
			ErrorColor:  cli.UiColorRed,
			Ui: &cli.BasicUi{
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
				Reader:      os.Stdin,
			},
		}}

	return RunCustom(args, Commands(meta))
}

func RunCustom(args []string, commands map[string]cli.CommandFactory) int {

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:       args,
		Commands:   commands,
		Version:    Version,
		HelpFunc:   cli.BasicHelpFunc(Name),
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
	}

	return exitCode
}
