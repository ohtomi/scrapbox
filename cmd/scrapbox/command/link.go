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

type LinkCommand struct {
	Meta
}

func (c *LinkCommand) FetchAllLinks(client *client.Client, project, page string) ([]string, error) {

	p, err := client.GetPage(context.Background(), project, page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page")
	}

	return p.ExtractExternalLinks(), nil
}

func (c *LinkCommand) Run(args []string) int {

	var (
		project string
		page    string

		token      string
		host       string
		expiration int
		userAgent  string
	)

	flags := flag.NewFlagSet("open", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&token, "token", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&host, "host", os.Getenv(EnvScrapboxHost), "")
	flags.StringVar(&host, "h", os.Getenv(EnvScrapboxHost), "")
	flags.IntVar(&expiration, "expire", EnvToInt(EnvExpiration, client.DefaultExpiration), "")
	flags.StringVar(&userAgent, "ua", os.Getenv(EnvUserAgent), "")

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

	if len(userAgent) == 0 {
		userAgent = client.DefaultUserAgent
	}

	// process

	client, err := client.NewClient(parsedURL, token, expiration, userAgent)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to initialize api client. cause: %s", err))
		return int(ExitCodeError)
	}

	linkURLs, err := c.FetchAllLinks(client, project, page)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to fetch the scrapbox page. cause: %s", err))
		return int(ExitCodeFetchFailure)
	}

	for _, u := range linkURLs {
		c.Ui.Output(u)
	}

	return int(ExitCodeOK)
}

func (c *LinkCommand) Synopsis() string {
	return "Print all URLs in the scrapbox page"
}

func (c *LinkCommand) Help() string {
	helpText := `usage: scrapbox link [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
  --host, -h   Scrapbox Host. By default, "https://scrapbox.io".
  --expire     Local Cache Expiration. By default, 3600 seconds.
  --ua         User Agent. By default, "ScrapboxGoClient/x.x.x"
`
	return strings.TrimSpace(helpText)
}
