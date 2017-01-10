package command

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	defaultURL = "https://scrapbox.io"
)

const (
	apiEndpoint = "api/pages"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxURL   = "SCRAPBOX_URL"
)

const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodeInvalidURL
	ExitCodeProjectNotFound
	ExitCodeTagNotFound
	ExitCodePageNotFound
	ExitCodeFetchFailure
	ExitCodeWriteRelatedPageFailure
	ExitCodeReadCacheFailure
	ExitCodeWriteCacheFailure
	ExitCodeOpenURLFailure
	ExitCodeNoAvailableURLFound
)

var (
	scrapboxHome = path.Join(os.Getenv("HOME"), ".scrapbox")
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui cli.Ui
}

func (m *Meta) TrimPortFromHost(host string) string {
	if strings.Index(host, ":") == -1 {
		return host
	} else {
		return host[:strings.Index(host, ":")]
	}
}
