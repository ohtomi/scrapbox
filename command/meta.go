package command

import (
	"os"
	"path"
	"strconv"

	"github.com/mitchellh/cli"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxHost  = "SCRAPBOX_HOST"
	EnvExpiration    = "SCRAPBOX_EXPIRATION"
)

type ExitCode int

const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodeInvalidURL
	ExitCodeProjectNotFound
	ExitCodePageNotFound
	ExitCodeFetchFailure
)

const (
	EnvHome        = "SCRAPBOX_HOME"
	EnvDebug       = "SCRAPBOX_DEBUG"
	EnvLongRunTest = "SCRAPBOX_LONG_RUN_TEST"
)

var (
	DebugMode       = os.Getenv(EnvDebug) != ""
	LongRunTestMode = os.Getenv(EnvLongRunTest) != ""
)

func ScrapboxHomeFromEnv() string {
	value := os.Getenv(EnvHome)
	if len(value) == 0 {
		value = path.Join(os.Getenv("HOME"), ".scrapbox")
	}
	return value
}

func EnvToInt(name string, value int) int {
	parsedInt, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return value
	}
	return parsedInt
}

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui cli.Ui
}
