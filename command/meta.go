package command

import "github.com/mitchellh/cli"

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui cli.Ui
}

const (
	defaultURL = "https://scrapbox.io"
)

const (
	apiEndpoint = "api/pages"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxURL   = "SCRAPBOX_URL"
)
