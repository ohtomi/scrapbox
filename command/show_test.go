package command

// scrapbox read --token xxx project page

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestShowCommand__show_english(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ShowCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "english"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
