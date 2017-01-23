package command

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	defaultURL  = "https://scrapbox.io"
	defaultHost = "scrapbox.io"

	defaultDownloadDir = "./"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxURL   = "SCRAPBOX_URL"
	EnvScrapboxHost  = "SCRAPBOX_HOST"

	EnvDownloadDir = "SCRAPBOX_DOWNLOAD_DIR"
)

const (
	apiEndPoint = "api/pages"
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

const (
	EnvHome  = "SCRAPBOX_HOME"
	EnvDebug = "SCRAPBOX_DEBUG"
)

var scrapboxHome string
var debugMode bool

func InitializeMeta() {

	scrapboxHome = os.Getenv(EnvHome)
	if len(scrapboxHome) == 0 {
		scrapboxHome = path.Join(os.Getenv("HOME"), ".scrapbox")
	}

	debugMode = os.Getenv(EnvDebug) != ""
}

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
