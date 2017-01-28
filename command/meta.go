package command

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	defaultHost = "https://scrapbox.io"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxHost  = "SCRAPBOX_HOST"
)

const (
	apiEndPoint = "api/pages"
)

type ExitCode int

const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodeInvalidURL
	ExitCodeProjectNotFound
	ExitCodeTagNotFound
	ExitCodePageNotFound
	ExitCodeFetchFailure
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
