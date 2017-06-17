package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type OpenCommand struct {
	Meta
}

func (c *OpenCommand) BuildPageURL(host, project, page string) string {
	return fmt.Sprintf("%s/%s/%s", host, project, EncodeURIComponent(page))
}

func (c *OpenCommand) Run(args []string) int {

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
		host = DefaultHost
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to parse the url. host: %s, cause: %s", host, err))
		return int(ExitCodeInvalidURL)
	}

	// process

	pageURL := c.BuildPageURL(host, project, page)
	c.Ui.Output(pageURL)

	return int(ExitCodeOK)
}

func (c *OpenCommand) Synopsis() string {
	return "Print the URL of the scrapbox page"
}

func (c *OpenCommand) Help() string {
	helpText := `usage: scrapbox open [options...] PROJECT PAGE

Options:
  --host, -h   Scrapbox Host. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
