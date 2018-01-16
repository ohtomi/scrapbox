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

type ListCommand struct {
	Meta
}

func (c *ListCommand) FetchRelatedPages(client *client.Client, project string, tags []string) ([]string, error) {

	q, err := client.ExecQuery(context.Background(), project, tags, 0, 100)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return q.Pages, nil
}

func (c *ListCommand) Run(args []string) int {

	var (
		project string
		tags    []string

		token      string
		host       string
		expiration int
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
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
		host = client.DefaultHost
	}

	parsedURL, err := url.ParseRequestURI(host)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to parse the url. host: %s, cause: %s", host, err))
		return int(ExitCodeInvalidURL)
	}

	// process

	client, err := client.NewClient(parsedURL, token, expiration)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to initialize api client. cause: %s", err))
		return int(ExitCodeError)
	}

	relatedPages, err := c.FetchRelatedPages(client, project, tags)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to fetch the scrapbox page. cause: %s", err))
		return int(ExitCodeFetchFailure)
	}

	for _, p := range relatedPages {
		c.Ui.Output(p)
	}

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
  --expire     Local Cache Expiration. By default, 3600 seconds.
`
	return strings.TrimSpace(helpText)
}
