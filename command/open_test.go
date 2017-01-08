package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestOpenCommand_implement(t *testing.T) {
	var _ cli.Command = &OpenCommand{}
}
