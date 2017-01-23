package command

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type DownloadCommand struct {
	Meta
}

func writeLocalFile(directory, page string, lines []string) error {

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

func (c *DownloadCommand) Run(args []string) int {

	var (
		project string
		page    string

		token   string
		baseURL string

		directory string

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
	flags.StringVar(&directory, "dest", os.Getenv(EnvDownloadDir), "")
	flags.StringVar(&directory, "d", os.Getenv(EnvDownloadDir), "")

	if err := flags.Parse(args); err != nil {
		return int(ExitCodeParseFlagsError)
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and PAGE name.")
		return int(ExitCodeBadArgs)
	}
	project, page = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return int(ExitCodeProjectNotFound)
	}
	if len(page) == 0 {
		c.Ui.Error("missing PAGE name.")
		return int(ExitCodePageNotFound)
	}

	if len(baseURL) == 0 {
		baseURL = defaultURL
	}
	if len(directory) == 0 {
		directory = defaultDownloadDir
	}

	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		c.Ui.Error("failed to parse url: " + baseURL)
		return int(ExitCodeInvalidURL)
	}
	host = c.Meta.TrimPortFromHost(parsedURL.Host)

	// process

	lines, err := fetchPageContent(host, project, page, token, parsedURL)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to fetch page: %s", err))
		return int(ExitCodeFetchFailure)
	}

	if err := writeLocalFile(directory, page, lines); err != nil {
		c.Ui.Error(fmt.Sprintf("failed to write local file: %s", err))
		return int(ExitCodeWriteFileFailure)
	}

	return int(ExitCodeOK)
}

func (c *DownloadCommand) Synopsis() string {
	return "Download remote content from scrapbox.io"
}

func (c *DownloadCommand) Help() string {
	helpText := `usage: scrapbox download [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid which is used to access private project.
  --url, -u    Scrapbox URL. By default, "https://scrapbox.io".
  --dest, -d   Output directory. By default, "./".
`
	return strings.TrimSpace(helpText)
}
