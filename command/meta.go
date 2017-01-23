package command

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	defaultURL     = "https://scrapbox.io"
	defaultHost    = "scrapbox.io"
	defaultDestDir = "./"
)

const (
	apiEndpoint = "api/pages"
)

const (
	EnvScrapboxToken   = "SCRAPBOX_TOKEN"
	EnvScrapboxURL     = "SCRAPBOX_URL"
	EnvScrapboxHost    = "SCRAPBOX_HOST"
	EnvScrapboxDestDir = "SCRAPBOX_DEST_DIR"

	EnvDebug = "SCRAPBOX_DEBUG"
)

const (
	ExitCodeOK int = 0
)

const (
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
	ExitCodeNoRelatedPagesFound
	ExitCodeOpenURLFailure
	ExitCodeNoAvailableURLFound
	ExitCodeWriteFileFailure
)

var (
	scrapboxHome = path.Join(os.Getenv("HOME"), ".scrapbox")
	debugMode    = false
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
