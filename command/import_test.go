package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestImportCommand_implement(t *testing.T) {
	var _ cli.Command = &ImportCommand{}
}
