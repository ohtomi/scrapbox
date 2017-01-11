package command

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type OpenCommand struct {
	Meta
}

func queryFirstURL(host, project, tag, relatedPage string) (string, error) {

	var firstURL = ""

	statement := "select first_url from related_page where host = ? and project = ? and lower(tag) = lower(?) and lower(related_page) = lower(?);"
	parameters := []interface{}{host, project, tag, relatedPage}
	handler := func(rows *sql.Rows) error {
		for rows.Next() {
			if err := rows.Scan(&firstURL); err != nil {
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
		return "", err
	}

	return firstURL, nil
}

func openUrl(firstURL string) error {

	if err := exec.Command("open", firstURL).Start(); err != nil {
		return err
	}
	return nil
}

func (c *OpenCommand) Run(args []string) int {

	var (
		project string
		tag     string
		page    string

		host string
	)

	flags := flag.NewFlagSet("open", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&host, "host", os.Getenv(EnvScrapboxHost), "")
	flags.StringVar(&host, "h", os.Getenv(EnvScrapboxHost), "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeParseFlagsError
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 3 {
		c.Ui.Error("you must set PROJECT, TAG, PAGE name.")
		return ExitCodeBadArgs
	}
	project, tag, page = parsedArgs[0], parsedArgs[1], parsedArgs[2]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return ExitCodeProjectNotFound
	}
	if len(tag) == 0 {
		c.Ui.Error("missing TAG name.")
		return ExitCodeTagNotFound
	}
	if len(page) == 0 {
		c.Ui.Error("missing PAGE name.")
		return ExitCodePageNotFound
	}

	if len(host) == 0 {
		host = defaultHost
	}

	// process

	firstURL, err := queryFirstURL(host, project, tag, page)
	if err != nil || len(firstURL) == 0 {
		c.Ui.Warn("no available url found.")
		return ExitCodeNoAvailableURLFound
	}

	if err := openUrl(firstURL); err != nil {
		c.Ui.Error(fmt.Sprintf("failed to open url: %s", err))
		return ExitCodeOpenURLFailure
	}
	return ExitCodeOK
}

func (c *OpenCommand) Synopsis() string {
	return "Open the first URL, which is embedded in specified page"
}

func (c *OpenCommand) Help() string {
	helpText := `usage: scrapbox open [options...] HOST PROJECT TAG PAGE

Options:
  --host, -h   Scrapbox Host. By default, "scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
