package command

import (
	"strings"
)

type OpenCommand struct {
	Meta
}

func (c *OpenCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *OpenCommand) Synopsis() string {
	return "Open the first URL, which is embedded in specified page"
}

func (c *OpenCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
