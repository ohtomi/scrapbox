package command

import (
	"strings"
)

type ShowCommand struct {
	Meta
}

func (c *ShowCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *ShowCommand) Synopsis() string {
	return "Show page content, the title of which is equal to specified"
}

func (c *ShowCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
