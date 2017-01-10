package command

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func countRelatedPage(host, project, tag string) (int, error) {

	var count int = 0

	statement := "select count(*) from related_page where host = ? and project = ? and page = ?;"
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

func TestImportCommand_implement(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := newTestMeta(outStream, errStream, inStream)
	command := &ImportCommand{
		Meta: *meta,
	}

	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)
	defer testAPIServer.Close()

	muxAPI.HandleFunc("/api/pages/ohtomi/Bookmark", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../testdata/Bookmark.json")
	})
	muxAPI.HandleFunc("/api/pages/ohtomi/GolangでAPI Clientを実装する | SOTA", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../testdata/GolangでAPI Clientを実装する | SOTA.json")
	})

	args := strings.Split("--url "+testAPIServer.URL+" ohtomi Bookmark", " ")
	exitStatus := command.Run(args)
	if exitStatus != ExitCodeOK {
		t.Fatalf("ExitStatus=%d, but want %d", exitStatus, ExitCodeOK)
	}

	actual, err := countRelatedPage("127.0.0.1", "ohtomi", "Bookmark")
	if err != nil {
		t.Fatalf("failed to count related page: %s", err)
	}

	expected := 1
	if actual != expected {
		t.Fatalf("Output=%d, but want %d", actual, expected)
	}
}
