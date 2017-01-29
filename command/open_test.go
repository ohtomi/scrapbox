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

func TestOpenCommand__print_url_having_paren(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having paren ( ) mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20paren%20(%20)%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_plus(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having plus + mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20plus%20%2B%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_question(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having question ? mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20question%20%3F%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_slash(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having slash / mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20slash%20%2F%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_whitespace(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having whitespaces", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20whitespaces"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_japanese(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  日本語タイトルのページ", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/%E6%97%A5%E6%9C%AC%E8%AA%9E%E3%82%BF%E3%82%A4%E3%83%88%E3%83%AB%E3%81%AE%E3%83%9A%E3%83%BC%E3%82%B8"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
