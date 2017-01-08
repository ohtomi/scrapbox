package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestDownloadCommand_implement(t *testing.T) {
	var _ cli.Command = &DownloadCommand{}
}
