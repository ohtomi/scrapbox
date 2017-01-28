package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {

	var (
		project string
		tags    []string

		token string
		host  string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
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
	if len(parsedArgs) < 1 {
		c.Ui.Error("you must set PROJECT.")
		return int(ExitCodeBadArgs)
	}
	project, tags = parsedArgs[0], parsedArgs[1:]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT.")
		return int(ExitCodeProjectNotFound)
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
	c.Ui.Info(fmt.Sprintf("%s %s %s %s", project, tags, token, host))

	return int(ExitCodeOK)
}

func (c *ListCommand) Synopsis() string {
	return "List page titles containing specified tags"
}

func (c *ListCommand) Help() string {
	helpText := `usage: scrapbox list [options...] PROJECT [TAGs...]

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
  --host, -h   Scrapbox Host. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
