package command

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type ImportCommand struct {
	Meta
}

func fetchRelatedPageTitlesByTag(host, project, tag string, client *Client) ([]string, error) {

	p, err := client.GetPage(context.Background(), project, tag)
	if err != nil {
		return nil, err
	}

	return p.RelatedPageTitles, nil
}

func fetchTagListAndFirstURLByTitle(host, project, relatedPage string, client *Client) (string, string, error) {

	p, err := client.GetPage(context.Background(), project, relatedPage)
	if err != nil {
		return "", "", err
	}

	return p.TagList(), p.FirstURL(), nil
}

func writeRelatedPage(host, project, tag, relatedPage, tagList, firstURL string) error {

	statement := "insert into related_page(host, project, tag, related_page, tag_list, first_url) values(?, ?, ?, ?, ?, ?);"
	parameters := []interface{}{host, project, tag, relatedPage, tagList, firstURL}
	if err := execSQL(statement, parameters); err != nil {
		return err
	}

	return nil
}

func clearRelatedPage(host, project, tag string) error {

	statement := "delete from related_page where host = ? and project = ? and tag = ?;"
	parameters := []interface{}{host, project, tag}
	if err := execSQL(statement, parameters); err != nil {
		return err
	}

	return nil
}

func (c *ImportCommand) Run(args []string) int {

	var (
		project string
		tag     string

		token   string
		baseURL string

		host string
	)

	flags := flag.NewFlagSet("import", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&token, "token", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvScrapboxToken), "")
	flags.StringVar(&baseURL, "url", os.Getenv(EnvScrapboxURL), "")
	flags.StringVar(&baseURL, "u", os.Getenv(EnvScrapboxURL), "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeParseFlagsError
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and TAG name.")
		return ExitCodeBadArgs
	}
	project, tag = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return ExitCodeProjectNotFound
	}
	if len(tag) == 0 {
		c.Ui.Error("missing TAG name.")
		return ExitCodeTagNotFound
	}

	if len(baseURL) == 0 {
		baseURL = defaultURL
	}

	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		c.Ui.Error("failed to parse url: " + baseURL)
		return ExitCodeInvalidURL
	}
	host = c.Meta.TrimPortFromHost(parsedURL.Host)

	// process

	client, err := NewClient(parsedURL, token)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to initialize api client: %s", err))
		return ExitCodeError
	}

	titles, err := fetchRelatedPageTitlesByTag(host, project, tag, client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to fetch related page titles: %s", err))
		return ExitCodeFetchFailure
	}

	if err := clearRelatedPage(host, project, tag); err != nil {
		c.Ui.Error(fmt.Sprintf("failed to delete related page: %s", err))
		return ExitCodeWriteRelatedPageFailure
	}

	for _, t := range titles {
		tagList, firstURL, err := fetchTagListAndFirstURLByTitle(host, project, t, client)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("failed to fetch tag list: %s", err))
			continue
		}
		if err := writeRelatedPage(host, project, tag, t, tagList, firstURL); err != nil {
			c.Ui.Error(fmt.Sprintf("failed to write related page: %s", err))
			return ExitCodeWriteRelatedPageFailure
		}
	}

	return ExitCodeOK
}

func (c *ImportCommand) Synopsis() string {
	return "Import contents from scrapbox.io to local cache database"
}

func (c *ImportCommand) Help() string {
	helpText := `usage: scrapbox import [options...] PROJECT TAG

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
