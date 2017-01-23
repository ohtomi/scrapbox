package command

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const (
	cacheExpiration = 60 * 60 * 24
)

type ShowCommand struct {
	Meta
}

func hasValidLocalCache(host, project, page string) bool {

	var cachedAt int64 = 0

	statement := "select cached_at from local_cache where host = ? and project = ? and page = ?;"
	parameters := []interface{}{host, project, page}
	handler := func(rows *sql.Rows) error {
		for rows.Next() {
			if err := rows.Scan(&cachedAt); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}

	err := querySQL(statement, parameters, handler)
	if err != nil {
		return false
	}

	diff := time.Now().Unix() - cachedAt
	return diff <= int64(cacheExpiration)
}

func canonicalFilepath(directory, filename string) string {
	escapedFilename := strings.Replace(filename, "/", "%2F", -1)
	return path.Join(directory, escapedFilename)
}

func readLocalCache(host, project, page string) ([]string, error) {

	var lines []string

	directory := path.Join(scrapboxHome, "page", host, project)
	filepath := canonicalFilepath(directory, page)

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

	statement := "insert or replace into local_cache(host, project, page, cached_at) values(?, ?, ?, ?);"
	parameters := []interface{}{host, project, page, time.Now().Unix()}
	if err := execSQL(statement, parameters); err != nil {
		return err
	}

	directory := path.Join(scrapboxHome, "page", host, project)
	filepath := canonicalFilepath(directory, page)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}
	fout, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fout.Close()

	for _, line := range lines {
		fout.WriteString(fmt.Sprintf("%s\n", line))
	}

	return nil
}

func fetchPageContent(host, project, page, token string, parsedURL *url.URL) ([]string, error) {

	client, err := NewClient(parsedURL, token)
	if err != nil {
		return nil, err
	}

	p, err := client.GetPage(context.Background(), project, page)
	if err != nil {
		return nil, err
	}

	return p.Lines, nil
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
	flags.BoolVar(&debugMode, "debug", false, "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeParseFlagsError
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and PAGE name.")
		return ExitCodeBadArgs
	}
	project, page = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return ExitCodeProjectNotFound
	}
	if len(page) == 0 {
		c.Ui.Error("missing PAGE name.")
		return ExitCodePageNotFound
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

	if !hasValidLocalCache(host, project, page) {
		lines, err := fetchPageContent(host, project, page, token, parsedURL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("failed to fetch page: %s", err))
			return ExitCodeFetchFailure
		}
		if err := writeLocalCache(host, project, page, lines); err != nil {
			c.Ui.Error(fmt.Sprintf("failed to write local cache: %s", err))
			return ExitCodeWriteCacheFailure
		}
	}

	lines, err := readLocalCache(host, project, page)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to read local cache: %s", err))
		return ExitCodeReadCacheFailure
	}
	for _, line := range lines {
		c.Ui.Output(line)
	}
	return ExitCodeOK
}

func (c *ShowCommand) Synopsis() string {
	return "Show page content, the title of which is equal to specified"
}

func (c *ShowCommand) Help() string {
	helpText := `usage: scrapbox show [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
