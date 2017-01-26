package command

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ListCommand struct {
	Meta
}

type RelatedPage struct {
	title, tagList string
}

func queryRelatedPages(host, project, tag string) ([]RelatedPage, error) {

	var relatedPages = []RelatedPage{}

	statement := "select related_page, tag_list from related_page where host = ? and project = ? and lower(tag) = lower(?);"
	parameters := []interface{}{host, project, tag}
	handler := func(rows *sql.Rows) error {
		for rows.Next() {
			var title, tagList string
			if err := rows.Scan(&title, &tagList); err != nil {
				return err
			}
			relatedPages = append(relatedPages, RelatedPage{title, tagList})
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}

	err := querySQL(statement, parameters, handler)
	if err != nil {
		return []RelatedPage{}, err
	}

	return relatedPages, nil
}

func (c *ListCommand) Run(args []string) int {

	var (
		project string
		tag     string

		host string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}

	flags.StringVar(&host, "host", os.Getenv(EnvScrapboxHost), "")
	flags.StringVar(&host, "h", os.Getenv(EnvScrapboxHost), "")

	if err := flags.Parse(args); err != nil {
		return int(ExitCodeParseFlagsError)
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 2 {
		c.Ui.Error("you must set PROJECT and TAG name.")
		return int(ExitCodeBadArgs)
	}
	project, tag = parsedArgs[0], parsedArgs[1]

	if len(project) == 0 {
		c.Ui.Error("missing PROJECT name.")
		return int(ExitCodeProjectNotFound)
	}
	if len(tag) == 0 {
		c.Ui.Error("missing TAG name.")
		return int(ExitCodeTagNotFound)
	}

	if len(host) == 0 {
		host = defaultHost
	}

	// process

	relatedPages, err := queryRelatedPages(host, project, tag)
	if err != nil || len(relatedPages) == 0 {
		return int(ExitCodeNoRelatedPagesFound)
	}

	for _, r := range relatedPages {
		c.Ui.Output(fmt.Sprintf("%s --- %s", r.title, r.tagList))
	}

	return int(ExitCodeOK)
}

func (c *ListCommand) Synopsis() string {
	return "List page titles include specified tag"
}

func (c *ListCommand) Help() string {
	helpText := `usage: scrapbox list [options...] PROJECT TAG

Options:
  --host, -h   Scrapbox Host. By default, "scrapbox.io".
`
	return strings.TrimSpace(helpText)
}
