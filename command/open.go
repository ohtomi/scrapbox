package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type OpenCommand struct {
	Meta
}

func containsUrl(line string) (bool, string) {

	keywords := []string{"http://", "https://"}
	whitespace := " "

	for _, k := range keywords {
		if strings.Contains(line, k) {
			if strings.Index(line, k) != -1 {
				line = line[strings.Index(line, k):]
			}
			if strings.Index(line, whitespace) != -1 {
				line = line[:strings.Index(line, whitespace)]
			}
			return true, line
		}
	}

	return false, ""
}

func openUrl(parsedURL string) error {

	if err := exec.Command("open", parsedURL).Start(); err != nil {
		return err
	}
	return nil
}

func (c *OpenCommand) Run(args []string) int {

	var (
		project string
		page    string

		token   string
		baseURL string

		host string
	)

	flags := flag.NewFlagSet("open", flag.ContinueOnError)
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
		if has, clippedURL := containsUrl(line); has {
			c.Ui.Output(clippedURL)
			if err := openUrl(clippedURL); err != nil {
				c.Ui.Error(fmt.Sprintf("failed to open url: %s", err))
				return ExitCodeOpenURLFailure
			}
			return ExitCodeOK
		}
	}
	c.Ui.Warn("no available url found.")
	return ExitCodeNoAvailableURLFound
}

func (c *OpenCommand) Synopsis() string {
	return "Open the first URL, which is embedded in specified page"
}

func (c *OpenCommand) Help() string {
	helpText := `usage: scrapbox open [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
