package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestShowCommand_implement(t *testing.T) {
	var _ cli.Command = &ShowCommand{}
}
