package command

import (
	"bytes"
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
		tag     string

		baseURL string

		host   string
		result bytes.Buffer
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
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

	fmt.Fprintf(&result, `Go Advent Calendar 2016 - Qiita --- #Go #adventcalendar #Qiita #Bookmark
Go (その2) Advent Calendar 2016 - Qiita --- #Go #adventcalendar #Qiita #Bookmark
Go (その3) Advent Calendar 2016 - Qiita --- #Go #adventcalendar #Qiita #Bookmark
高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA --- #gcli #Go #generator #Bookmark
...
`)
	c.Ui.Output(result.String())

	c.Ui.Output("debug: " + project + " " + tag + " " + baseURL + " " + host)
	return 0
}

func (c *ListCommand) Synopsis() string {
	return "List page titles include specified tag"
}

func (c *ListCommand) Help() string {
	helpText := `usage: scrapbox list [options...] PROJECT TAG

Options:
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
