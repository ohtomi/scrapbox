package command

import (
	"flag"
	"net/url"
	"os"
	"strings"
)

type DownloadCommand struct {
	Meta
}

func (c *DownloadCommand) Run(args []string) int {

	var (
		project string
		tag     string

		token   string
		baseURL string

		host string
	)

	flags := flag.NewFlagSet("download", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&token, "token", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&baseURL, "url", os.Getenv(EnvScrapboxURL), "")
	flags.StringVar(&baseURL, "u", os.Getenv(EnvScrapboxURL), "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and TAG name.")
		return 1
	}
	project, tag = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return 1
	}
	if len(tag) == 0 {
		c.Ui.Error("missing TAG name.")
		return 1
	}

	if len(baseURL) == 0 {
		baseURL = defaultURL
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		c.Ui.Error("failed to parse url: " + baseURL)
		return 1
	}
	host = u.Host

	c.Ui.Output("debug: " + project + " " + tag + " " + token + " " + baseURL + " " + host)
	return 0
}

func (c *DownloadCommand) Synopsis() string {
	return "Download remote content from scrapbox.io"
}

func (c *DownloadCommand) Help() string {
	helpText := `usage: scrapbox download [options...] PROJECT TAG

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
