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

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ShowCommand{
		Meta: *meta,
	}

	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)
	defer testAPIServer.Close()

	muxAPI.HandleFunc("/api/pages/go-scrapbox/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "/api/pages/go-scrapbox/", "", -1)
		escaped := strings.Replace(path, "/", "%2F", -1)
		http.ServeFile(w, r, fmt.Sprintf("../testdata/scrapbox.io/go-scrapbox/%s", escaped))
	})

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
