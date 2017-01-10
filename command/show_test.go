package command

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestShowCommand_implement(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ShowCommand{
		Meta: *meta,
	}

	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)
	defer testAPIServer.Close()

	muxAPI.HandleFunc("/api/pages/Bookmark", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../testdata/Bookmark.json")
	})

	args := strings.Split("--url "+testAPIServer.URL+" ohtomi Bookmark", " ")
	exitStatus := command.Run(args)
	if exitStatus != ExitCodeOK {
		t.Fatalf("ExitStatus=%d, but want %d", exitStatus, ExitCodeOK)
	}

	expected := fmt.Sprintf("Bookmark")
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output=%q, but want %q", outStream.String(), expected)
	}
}
