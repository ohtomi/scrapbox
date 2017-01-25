package command

import (
	"bytes"
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func countRelatedPage(host, project, tag string) (int, error) {

	var count int = 0

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
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ImportCommand{
		Meta: *meta,
	}

	testAPIServer := runAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox english", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus actual %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	actual, err := countRelatedPage("127.0.0.1", "go-scrapbox", "english")
	if err != nil {
		t.Fatalf("failed to count related page: %s", err)
	}

	expected := 5
	if actual != expected {
		t.Fatalf("Count actual %d, but want %d", actual, expected)
	}
}

func TestImportCommand__import_japanese_pages(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ImportCommand{
		Meta: *meta,
	}

	testAPIServer := runAPIServer()
	defer testAPIServer.Close()

	args := strings.Split("--url "+testAPIServer.URL+" go-scrapbox japanese", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus actual %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	actual, err := countRelatedPage("127.0.0.1", "go-scrapbox", "japanese")
	if err != nil {
		t.Fatalf("failed to count related page: %s", err)
	}

	expected := 6
	if actual != expected {
		t.Fatalf("Count actual %d, but want %d", actual, expected)
	}
}
