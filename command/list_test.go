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

	"github.com/MakeNowJust/heredoc"
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
		tags := strings.Split(query.Get("q"), " ")

		filename := fmt.Sprintf("%s-%s", skip, limit)
		directory := path.Join("../testdata/query/scrapbox.io/go-scrapbox", path.Join(tags...))
		filepath := path.Join(directory, EncodeFilename(filename))
		http.ServeFile(w, r, filepath)
	})

	muxAPI.HandleFunc("/api/pages/go-scrapbox", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		skip := query.Get("skip")
		limit := query.Get("limit")

		filename := fmt.Sprintf("%s-%s", skip, limit)
		directory := path.Join("../testdata/query/scrapbox.io/go-scrapbox")
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

func TestListCommand__find_by_project_only(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ListCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--host "+testAPIServer.URL+" go-scrapbox", " ")
	exitStatus := command.Run(args)

	if DebugMode {
		t.Log(outStream.String())
		t.Log(errStream.String())
	}

	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := heredoc.Doc(`
		HTTPなリンクのあるページ
		HTTPSなリンクのあるページ
		title having question ? mark
		title having plus + mark
		title having paren ( ) mark
		title having slash / mark
		文章のなかにリンクがあるページ2
		文章のなかにリンクがあるページ1
		複数のリンクがあるページ
		title having whitespaces
		日本語タイトルのページ
	`)
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is \n%s\n, but want \n%s", outStream.String(), expected)
	}
}

func TestListCommand__find_by_project_and_one_keyword(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ListCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--host "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)

	if DebugMode {
		t.Log(outStream.String())
		t.Log(errStream.String())
	}

	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := heredoc.Doc(`
		title having question ? mark
		title having plus + mark
		title having paren ( ) mark
		title having slash / mark
		title having whitespaces
	`)
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is \n%s\n, but want \n%s", outStream.String(), expected)
	}
}

func TestListCommand__find_by_project_and_many_keywords(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ListCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--host "+testAPIServer.URL+" go-scrapbox english paren", " ")
	exitStatus := command.Run(args)

	if DebugMode {
		t.Log(outStream.String())
		t.Log(errStream.String())
	}

	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "title having paren ( ) mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestListCommand__find_by_project_and_non_tag_keyword(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ListCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--host "+testAPIServer.URL+" go-scrapbox english whitespaces", " ")
	exitStatus := command.Run(args)

	if DebugMode {
		t.Log(outStream.String())
		t.Log(errStream.String())
	}

	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "title having whitespaces"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
