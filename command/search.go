package command

import (
	"strings"
)

type SearchCommand struct {
	Meta
}

func (c *SearchCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *SearchCommand) Synopsis() string {
	return "Search by keyword"
}

func (c *SearchCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
