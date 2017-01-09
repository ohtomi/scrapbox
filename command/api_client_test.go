package command

import (
	"io"
	"os"
	_ "testing"

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
