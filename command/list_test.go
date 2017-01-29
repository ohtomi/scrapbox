package command

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestListCommand__find_by_english(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ListCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected :=
		`title having question ? mark
title having plus + mark
title having paren ( ) mark
title having slash / mark
title having whitespaces
`
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
