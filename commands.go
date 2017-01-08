package main

import (
	"github.com/mitchellh/cli"
	"github.com/ohtomi/scrapbox/command"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"search": func() (cli.Command, error) {
			return &command.SearchCommand{
				Meta: *meta,
			}, nil
		},
		"import": func() (cli.Command, error) {
			return &command.ImportCommand{
				Meta: *meta,
			}, nil
		},
		"upload": func() (cli.Command, error) {
			return &command.UploadCommand{
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
