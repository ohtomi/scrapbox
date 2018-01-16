package main

import (
	"github.com/mitchellh/cli"
	"github.com/ohtomi/scrapbox/cmd/scrapbox/command"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				Meta: *meta,
			}, nil
		},
		"open": func() (cli.Command, error) {
			return &command.OpenCommand{
				Meta: *meta,
			}, nil
		},
		"link": func() (cli.Command, error) {
			return &command.LinkCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
