package command

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type ShowCommand struct {
	Meta
}

func hasValidLocalCache(host, project, page string) bool {

	// TODO check file timestamp
	return false
}

func readLocalCache(host, project, page string) ([]string, error) {

	var lines []string

	// TODO local cache file - path, name
	filepath := os.Getenv("HOME") + "/.scrapbox/" + host + "_" + page
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func writeLocalCache(host, project, page string, lines []string) error {

	// TODO local cache file - path, name
	filepath := os.Getenv("HOME") + "/.scrapbox/" + host + "_" + page
	// TODO create directory if not exists
	fout, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fout.Close()

	for _, l := range lines {
		fout.WriteString(l)
	}

	return nil
}

func (c *ShowCommand) Run(args []string) int {

	var (
		project string
		page    string

		token   string
		baseURL string

		host string
	)

	flags := flag.NewFlagSet("show", flag.ContinueOnError)
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
		c.Ui.Error("you must set PROJECT and PAGE name.")
		return 1
	}
	project, page = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return 1
	}
	if len(page) == 0 {
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

	if hasValidLocalCache(host, project, page) {
		lines, err := readLocalCache(host, project, page)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("failed to read local cache: %s", err))
			return 1
		}
		for _, l := range lines {
			c.Ui.Output(l)
		}
		return 0
	}

	client, err := NewClient(u, token)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to initialize http client: %s", err))
	}

	p, err := client.GetPage(context.Background(), page)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get page: %s", err))
		return 1
	}

	err = writeLocalCache(host, project, page, p.Lines)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to write local cache: %s", err))
		return 1
	}

	for _, l := range p.Lines {
		c.Ui.Output(l)
	}

	return 0
}

func (c *ShowCommand) Synopsis() string {
	return "Show page content, the title of which is equal to specified"
}

func (c *ShowCommand) Help() string {
	helpText := `usage: scrapbox show [options...] PROJECT TAG

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
