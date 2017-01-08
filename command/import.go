package command

import (
	"strings"
)

type ImportCommand struct {
	Meta
}

func (c *ImportCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *ImportCommand) Synopsis() string {
	return "Import contents from scrapbox.io to local cache database"
}

func (c *ImportCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
