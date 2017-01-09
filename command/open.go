package command

import (
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type OpenCommand struct {
	Meta
}

func (c *OpenCommand) Run(args []string) int {

	var (
		project string
		tag     string

		baseURL string

		host   string
		result bytes.Buffer
	)

	flags := flag.NewFlagSet("open", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

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

	fmt.Fprintf(&result, "-> open http://deeeet.com/writing/2014/06/22/cli-init/")
	c.Ui.Output(result.String())

	c.Ui.Output("debug: " + project + " " + tag + " " + baseURL + " " + host)
	return 0
}

func (c *OpenCommand) Synopsis() string {
	return "Open the first URL, which is embedded in specified page"
}

func (c *OpenCommand) Help() string {
	helpText := `usage: scrapbox open [options...] PROJECT TAG

Options:
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
