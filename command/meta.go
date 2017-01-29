package command

import (
	"os"
	"path"

	"github.com/mitchellh/cli"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxHost  = "SCRAPBOX_HOST"
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

var ScrapboxHome string
var DebugMode bool

func InitializeMeta() {

	ScrapboxHome = os.Getenv(EnvHome)
	if len(ScrapboxHome) == 0 {
		ScrapboxHome = path.Join(os.Getenv("HOME"), ".scrapbox")
	}

	DebugMode = os.Getenv(EnvDebug) != ""
}

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui cli.Ui
}
