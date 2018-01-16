package command

import (
	"bytes"
	"fmt"
)

type VersionCommand struct {
	Meta

	Name     string
	Version  string
	Revision string
}

func (c *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s version %s", c.Name, c.Version)
	if c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	c.Ui.Output(versionString.String())
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

func (c *VersionCommand) Help() string {
	return ""
}
