package command

import (
	"os"
	"strconv"

	"github.com/mitchellh/cli"
)

const (
	EnvScrapboxToken = "SCRAPBOX_TOKEN"
	EnvScrapboxHost  = "SCRAPBOX_HOST"
	EnvExpiration    = "SCRAPBOX_EXPIRATION"
	EnvUserAgent     = "SCRAPBOX_USER_AGENT"
)

const (
	EnvDebug       = "SCRAPBOX_DEBUG"
	EnvLongRunTest = "SCRAPBOX_LONG_RUN_TEST"
)

var (
	DebugMode       = os.Getenv(EnvDebug) != ""
	LongRunTestMode = os.Getenv(EnvLongRunTest) != ""
)

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
