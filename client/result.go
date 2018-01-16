package client

import (
	"fmt"
	"regexp"
	"strings"
)

type QueryResult struct {
	Count int
	Pages []string
}

type Page struct {
	Title string
	Lines []string
	Links []string
}

func (p *Page) ExtractExternalLinks() []string {

	includes := []string{"http://", "https://"}
	excludes := []string{".png", ".gif", ".jpg", ".jpeg", ".svg"}
	whitespace := " "

	match := func(line string, keywords []string) string {
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				return keyword
			}
		}
		return ""
	}

	linkURLs := []string{}

	for _, line := range p.Lines {
		if matched := match(line, includes); matched != "" {
			if match(line, excludes) != "" {
				continue
			}
			foundBracket, _ := regexp.MatchString(fmt.Sprintf("\\[.*%s.*\\]", matched), line)
			if strings.Index(line, matched) != -1 {
				line = line[strings.Index(line, matched):]
			}
			if strings.Index(line, whitespace) != -1 {
				line = line[:strings.Index(line, whitespace)]
			}
			if foundBracket && strings.Index(line, "]") == len(line)-1 {
				line = line[:len(line)-1]
			}
			linkURLs = append(linkURLs, line)
		}
	}

	return linkURLs
}
