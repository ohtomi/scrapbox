package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type LinkCommand struct {
	Meta
}

func (c *LinkCommand) Run(args []string) int {

	var (
		project string
		page    string

		host string
	)

	flags := flag.NewFlagSet("open", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&host, "host", os.Getenv(EnvScrapboxHost), "")
	flags.StringVar(&host, "h", os.Getenv(EnvScrapboxHost), "")

	if err := flags.Parse(args); err != nil {
		return int(ExitCodeParseFlagsError)
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and PAGE.")
		return int(ExitCodeBadArgs)
	}
	project, page = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT.")
		return int(ExitCodeProjectNotFound)
	}
	if len(page) == 0 {
		c.Ui.Error("missing PAGE.")
		return int(ExitCodePageNotFound)
	}

	if len(host) == 0 {
		host = defaultHost
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		c.Ui.Error("failed to parse url: " + host)
		return int(ExitCodeInvalidURL)
	}

	// process
	c.Ui.Info(fmt.Sprintf("%s %s %s", project, page, host))

	return int(ExitCodeOK)
}

func (c *LinkCommand) Synopsis() string {
	return "Open each URLs, written in the scrapbox page, in the browser"
}

func (c *LinkCommand) Help() string {
	helpText := `usage: scrapbox link [options...] PROJECT PAGE

Options:
	--host, -h   Scrapbox Host. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
