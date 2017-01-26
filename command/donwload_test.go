package command

// scrapbox download --url http://xxx --token xxx project page /path/to/file

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestDownloadCommand__todo(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &DownloadCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" --dest ../testdata/download/ go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}
}
