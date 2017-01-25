package command

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

func setTestEnv(key, val string) func() {

	preVal := os.Getenv(key)
	os.Setenv(key, val)
	return func() {
		os.Setenv(key, preVal)
	}
}

func newTestMeta(outStream, errStream io.Writer, inStream io.Reader) *Meta {
	return &Meta{
		Ui: &cli.BasicUi{
			Writer:      outStream,
			ErrorWriter: errStream,
			Reader:      inStream,
		}}
}

func runAPIServer() *httptest.Server {

	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)

	muxAPI.HandleFunc("/api/pages/go-scrapbox/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "/api/pages/go-scrapbox/", "", -1)
		escaped := strings.Replace(path, "/", "%2F", -1)
		http.ServeFile(w, r, fmt.Sprintf("../testdata/scrapbox.io/go-scrapbox/%s", escaped))
	})

	return testAPIServer
}
