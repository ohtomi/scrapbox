package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type ReadCommand struct {
	Meta
}

func (c *ReadCommand) Run(args []string) int {

	var (
		project string
		page    string

		token string
		host  string
	)

	flags := flag.NewFlagSet("read", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&token, "token", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvScrapboxToken), "")
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
	c.Ui.Info(fmt.Sprintf("%s %s %s %s", project, page, token, host))

	return int(ExitCodeOK)
}

func (c *ReadCommand) Synopsis() string {
	return "Print content of the scrapbox page"
}

func (c *ReadCommand) Help() string {
	helpText := `usage: scrapbox read [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
	--host, -h   Scrapbox Host. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
