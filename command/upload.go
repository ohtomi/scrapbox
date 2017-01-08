package command

import (
	"strings"
)

type UploadCommand struct {
	Meta
}

func (c *UploadCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *UploadCommand) Synopsis() string {
	return "Upload local content to scrapbox.io"
}

func (c *UploadCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
