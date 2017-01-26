package command

// scrapbox import --url http://xxx --token xxx project tag

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

	muxAPI.HandleFunc("/api/pages/go-scrapbox/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "/api/pages/go-scrapbox/", "", -1)
		escaped := strings.Replace(path, "/", "%2F", -1)
		http.ServeFile(w, r, fmt.Sprintf("../testdata/scrapbox.io/go-scrapbox/%s", escaped))
	})

	return testAPIServer
}

func countRelatedPage(host, project, tag string) (int, error) {

	count := 0

	statement := "select count(*) from related_page where host = ? and project = ? and tag = ?;"
	parameters := []interface{}{host, project, tag}
	handler := func(rows *sql.Rows) error {
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}

	err := querySQL(statement, parameters, handler)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func TestImportCommand__import_english_pages(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ImportCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	actual, err := countRelatedPage("127.0.0.1", "go-scrapbox", "english")
	if err != nil {
		t.Fatalf("failed to count related page: %s", err)
	}

	expected := 5
	if actual != expected {
		t.Fatalf("Count is %d, but want %d", actual, expected)
	}
}

func TestImportCommand__import_japanese_pages(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &ImportCommand{
		Meta: *meta,
	}

	testAPIServer := RunAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox japanese", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	actual, err := countRelatedPage("127.0.0.1", "go-scrapbox", "japanese")
	if err != nil {
		t.Fatalf("failed to count related page: %s", err)
	}

	expected := 6
	if actual != expected {
		t.Fatalf("Count is %d, but want %d", actual, expected)
	}
}
