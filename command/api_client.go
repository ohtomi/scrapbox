package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultHost       = "https://scrapbox.io"
	DefaultExpiration = 60 * 60 * time.Second
)

const (
	userAgent = "ScrapboxGoClient/0.1.0"
)

func EncodeURIComponent(component string) string {
	regularEscaped := url.QueryEscape(component)
	rParenUnescaped := strings.Replace(regularEscaped, "%28", "(", -1)
	lParenUnescaped := strings.Replace(rParenUnescaped, "%29", ")", -1)
	plusEscaped := strings.Replace(lParenUnescaped, "+", "%20", -1)
	return plusEscaped
}

func buildQueryPath(project string, tags []string, skip, limit int) string {
	params := fmt.Sprintf("skip=%d&sort=updated&limit=%d&q=%s", skip, limit, EncodeURIComponent(strings.Join(tags, " ")))
	if len(tags) == 0 {
		return fmt.Sprintf("api/pages/%s?%s", project, params)
	} else {
		return fmt.Sprintf("api/pages/%s/search/query?%s", project, params)
	}
}

func buildPagePath(project, page string) string {
	return fmt.Sprintf("api/pages/%s/%s", project, EncodeURIComponent(page))
}

func trimPortFromHost(host string) string {
	if index := strings.Index(host, ":"); index == -1 {
		return host
	} else {
		return host[:index]
	}
}

func EncodeFilename(filename string) string {
	return strings.Replace(filename, "/", "%2F", -1)
}

func createQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(ScrapboxHome, "query", trimPortFromHost(host), project, path.Join(tags...))
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fout, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return fout, nil
}

func haveGoodQueryResultFile(host, project string, tags []string, skip, limit int, expiration time.Duration) bool {

	directory := path.Join(ScrapboxHome, "query", trimPortFromHost(host), project, path.Join(tags...))
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	if finfo.IsDir() {
		return false
	}
	mod := finfo.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(ScrapboxHome, "query", trimPortFromHost(host), project, path.Join(tags...))
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return fin, nil
}

func createPageFile(host, project, page string) (*os.File, error) {

	directory := path.Join(ScrapboxHome, "page", trimPortFromHost(host), project)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}
	filepath := path.Join(directory, EncodeFilename(page))
	fout, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return fout, nil
}

func haveGoodPageFile(host, project, page string, expiration time.Duration) bool {

	directory := path.Join(ScrapboxHome, "page", trimPortFromHost(host), project)
	filepath := path.Join(directory, EncodeFilename(page))
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	if finfo.IsDir() {
		return false
	}
	mod := finfo.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openPageFile(host, project, page string) (*os.File, error) {

	directory := path.Join(ScrapboxHome, "page", trimPortFromHost(host), project)
	filepath := path.Join(directory, EncodeFilename(page))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return fin, nil
}

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	Token string
}

func NewClient(url *url.URL, token string) (*Client, error) {
	// TODO proxy, ssl, timeout
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
		Token:      token,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {

	baseURL := *c.URL
	u := fmt.Sprintf("%s/%s", baseURL.String(), spath)

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	if len(c.Token) != 0 {
		req.Header.Set("Cookie", "connect.sid="+c.Token)
	}

	return req, nil
}

func (c *Client) decodeBody(resp *http.Response, out interface{}, f *os.File) error {
	defer resp.Body.Close()
	if f != nil {
		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
		defer f.Close()
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) decodeFromFile(resp *os.File, out interface{}) error {
	defer resp.Close()
	decoder := json.NewDecoder(resp)
	return decoder.Decode(out)
}

type QueryResult struct {
	Count int
	Pages []string
}

func (c *Client) ExecQuery(ctx context.Context, project string, tags []string, skip, limit int) (*QueryResult, error) {

	var (
		count int
		pages []string

		v interface{}
	)

	host := (*c.URL).Host
	expiration := DefaultExpiration // TODO
	if haveGoodQueryResultFile(host, project, tags, skip, limit, expiration) {
		res, err := openQueryResultFile(host, project, tags, skip, limit)
		if err != nil {
			return nil, err
		}
		if err := c.decodeFromFile(res, &v); err != nil {
			return nil, err
		}
	} else {
		spath := buildQueryPath(project, tags, skip, limit)
		req, err := c.newRequest(ctx, "GET", spath, nil)
		if err != nil {
			return nil, err
		}

		res, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Check status code here…
		if res.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("http status is %q", res.Status))
		}

		fout, err := createQueryResultFile(host, project, tags, skip, limit)
		if err != nil {
			return nil, err
		}

		if err := c.decodeBody(res, &v, fout); err != nil {
			return nil, err
		}
	}

	for _, p := range v.(interface{}).(map[string]interface{})["pages"].([]interface{}) {
		if len(tags) > 0 {
			for _, s := range p.(map[string]interface{})["snipet"].([]interface{}) {
				all := true
				for _, t := range tags {
					all = all &&
						(strings.Contains(strings.ToLower(s.(string)), fmt.Sprintf("<b>%s</b>", strings.ToLower(t))) ||
							strings.Contains(strings.ToLower(p.(map[string]interface{})["title"].(interface{}).(string)), strings.ToLower(t)))
				}
				if all {
					pages = append(pages, p.(map[string]interface{})["title"].(interface{}).(string))
					break
				}
			}
		} else {
			pages = append(pages, p.(map[string]interface{})["title"].(interface{}).(string))
		}
	}

	count = int(v.(interface{}).(map[string]interface{})["count"].(float64))
	if count > limit+skip || count == limit {
		q, err := c.ExecQuery(context.Background(), project, tags, skip+limit, limit)
		if err != nil {
			return nil, err
		}
		pages = append(pages, q.Pages...)
	}

	return &QueryResult{
		Count: count,
		Pages: pages,
	}, nil
}

type Page struct {
	Title string
	Lines []string
	Links []string
}

func (c *Client) GetPage(ctx context.Context, project, page string) (*Page, error) {

	var (
		v interface{}
	)

	host := (*c.URL).Host
	expiration := DefaultExpiration // TODO
	if haveGoodPageFile(host, project, page, expiration) {
		res, err := openPageFile(host, project, page)
		if err != nil {
			return nil, err
		}
		if err := c.decodeFromFile(res, &v); err != nil {
			return nil, err
		}
	} else {
		spath := buildPagePath(project, page)
		req, err := c.newRequest(ctx, "GET", spath, nil)
		if err != nil {
			return nil, err
		}

		res, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Check status code here…
		if res.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("http status is %q", res.Status))
		}

		fout, err := createPageFile(host, project, page)
		if err != nil {
			return nil, err
		}

		if err := c.decodeBody(res, &v, fout); err != nil {
			return nil, err
		}
	}

	title := v.(interface{}).(map[string]interface{})["title"].(string)
	lines := make([]string, len(v.(interface{}).(map[string]interface{})["lines"].([]interface{})))
	for i, l := range v.(interface{}).(map[string]interface{})["lines"].([]interface{}) {
		lines[i] = l.(map[string]interface{})["text"].(interface{}).(string)
	}
	links := make([]string, len(v.(interface{}).(map[string]interface{})["links"].([]interface{})))
	for i, l := range v.(interface{}).(map[string]interface{})["links"].([]interface{}) {
		links[i] = l.(string)
	}

	return &Page{
		Title: title,
		Lines: lines,
		Links: links,
	}, nil
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
