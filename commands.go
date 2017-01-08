package main

import (
	"github.com/mitchellh/cli"
	"github.com/ohtomi/scrapbox/command"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"import": func() (cli.Command, error) {
			return &command.ImportCommand{
				Meta: *meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"show": func() (cli.Command, error) {
			return &command.ShowCommand{
				Meta: *meta,
			}, nil
		},
		"open": func() (cli.Command, error) {
			return &command.OpenCommand{
				Meta: *meta,
			}, nil
		},
		"download": func() (cli.Command, error) {
			return &command.DownloadCommand{
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
