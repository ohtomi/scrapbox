package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestSearchCommand_implement(t *testing.T) {
	var _ cli.Command = &SearchCommand{}
}
