package command

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ohtomi/scrapbox/client"
	"github.com/pkg/errors"
)

type ReadCommand struct {
	Meta
}

func (c *ReadCommand) FetchContent(client *client.Client, project, page string) ([]string, error) {

	p, err := client.GetPage(context.Background(), project, page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page")
	}

	return p.Lines, nil
}

func (c *ReadCommand) Run(args []string) int {

	var (
		project string
		page    string

		token      string
		host       string
		expiration int
	)

	flags := flag.NewFlagSet("read", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&token, "token", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&host, "host", os.Getenv(EnvScrapboxHost), "")
	flags.StringVar(&host, "h", os.Getenv(EnvScrapboxHost), "")
	flags.IntVar(&expiration, "expire", EnvToInt(EnvExpiration, client.DefaultExpiration), "")

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
		host = client.DefaultHost
	}

	parsedURL, err := url.ParseRequestURI(host)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to parse the url. host: %s, cause: %s", host, err))
		return int(ExitCodeInvalidURL)
	}

	// process

	client, err := client.NewClient(ScrapboxHomeFromEnv(), parsedURL, token, expiration)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to initialize api client. cause: %s", err))
		return int(ExitCodeError)
	}

	lines, err := c.FetchContent(client, project, page)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to fetch the scrapbox page. cause: %s", err))
		return int(ExitCodeFetchFailure)
	}

	for _, l := range lines {
		c.Ui.Output(l)
	}

	return int(ExitCodeOK)
}

func (c *ReadCommand) Synopsis() string {
	return "Print the content of the scrapbox page"
}

func (c *ReadCommand) Help() string {
	helpText := `usage: scrapbox read [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
  --host, -h   Scrapbox Host. By default, "https://scrapbox.io".
  --expire     Local Cache Expiration. By default, 3600 seconds.
`
	return strings.TrimSpace(helpText)
}
