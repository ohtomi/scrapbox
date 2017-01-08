package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestUploadCommand_implement(t *testing.T) {
	var _ cli.Command = &UploadCommand{}
}
