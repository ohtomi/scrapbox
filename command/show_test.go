package command

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestShowCommand__show_english(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ShowCommand{
		Meta: *meta,
	}

	testAPIServer := runAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus actual %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "english"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output actual %q, but want %q", outStream.String(), expected)
	}
}
