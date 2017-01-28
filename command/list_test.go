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

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--host "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

// 	expected :=
// 		`title having paren ( ) mark --- #english #no-url #whitespace #no-slash #paren #no-plus #no-question
// title having plus + mark --- #english #no-url #whitespace #no-slash #no-paren #plus #no-question
// title having question ? mark --- #english #no-url #whitespace #no-slash #no-paren #no-plus #question
// title having slash / mark --- #english #no-url #whitespace #slash #no-paren #no-plus #no-question
// title having whitespaces --- #english #no-url #whitespace #no-slash #no-paren #no-plus #no-question
// `
// 	if !strings.Contains(outStream.String(), expected) {
// 		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
// 	}
}
