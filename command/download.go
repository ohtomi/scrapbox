package command

import (
	"strings"
)

type DownloadCommand struct {
	Meta
}

func (c *DownloadCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *DownloadCommand) Synopsis() string {
	return "Download remote content from scrapbox.io"
}

func (c *DownloadCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
