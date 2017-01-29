package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func SetTestEnv(key, newValue string) func() {

	prevValue := os.Getenv(key)
	os.Setenv(key, newValue)
	reverter := func() {
		os.Setenv(key, prevValue)
	}
	return reverter
}

func NewTestMeta(outStream, errStream io.Writer, inStream io.Reader) *Meta {

	return &Meta{
		Ui: &cli.BasicUi{
			Writer:      outStream,
			ErrorWriter: errStream,
			Reader:      inStream,
		}}
}

func RunAPIServer() *httptest.Server {

	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)

	muxAPI.HandleFunc("/api/pages/go-scrapbox/search/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		skip := query.Get("skip")
		limit := query.Get("limit")
		tags := query["q"]

		filename := fmt.Sprintf("%s-%s", skip, limit)
		directory := path.Join("../testdata/query/scrapbox.io/go-scrapbox", path.Join(tags...))
		filepath := path.Join(directory, EncodeFilename(filename))
		http.ServeFile(w, r, filepath)
	})

	muxAPI.HandleFunc("/api/pages/go-scrapbox/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path

		filename := strings.Replace(urlPath, "/api/pages/go-scrapbox/", "", -1)
		directory := "../testdata/page/scrapbox.io/go-scrapbox"
		filepath := path.Join(directory, EncodeFilename(filename))
		http.ServeFile(w, r, filepath)
	})

	return testAPIServer
}

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
