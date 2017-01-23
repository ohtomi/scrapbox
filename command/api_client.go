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
	"strings"
)

const (
	userAgent = "ScrapboxGoClient/0.1.0"
)

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

func (c *Client) encodeURIComponent(component string) string {
	regularEscaped := url.QueryEscape(component)
	rParenUnescaped := strings.Replace(regularEscaped, "%28", "(", -1)
	lParenUnescaped := strings.Replace(rParenUnescaped, "%29", ")", -1)
	plusEscaped := strings.Replace(lParenUnescaped, "+", "%20", -1)
	return plusEscaped
}

type Page struct {
	Title             string
	Lines             []string
	Links             []string
	RelatedPageTitles []string
}

func (c *Client) GetPage(ctx context.Context, project, page string) (*Page, error) {
	spath := fmt.Sprintf("%s/%s/%s", apiEndPoint, project, c.encodeURIComponent(page))
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

	// Check debug mode
	var fout *os.File
	if debugMode {
		host := (*c.URL).Host
		directory := path.Join("testdata", host, project)
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return nil, err
		}
		filepath := canonicalFilepath(directory, page)
		fout, err = os.Create(filepath)
		if err != nil {
			return nil, err
		}
	}

	var v interface{}
	if err := c.decodeBody(res, &v, fout); err != nil {
		return nil, err
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
	relatedPageTitles := make([]string, len(v.(interface{}).(map[string]interface{})["relatedPages"].([]interface{})))
	for i, r := range v.(interface{}).(map[string]interface{})["relatedPages"].([]interface{}) {
		relatedPageTitles[i] = r.(map[string]interface{})["title"].(interface{}).(string)
	}

	return &Page{
		Title:             title,
		Lines:             lines,
		Links:             links,
		RelatedPageTitles: relatedPageTitles,
	}, nil
}

func (p *Page) TagList() string {
	var tagList = ""
	for _, l := range p.Links {
		tagList = fmt.Sprintf("%s #%s", tagList, l)
	}
	return strings.TrimSpace(tagList)
}

func (p *Page) FirstURL() string {

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

	for _, line := range p.Lines {
		if matched := match(line, includes); matched != "" {
			if match(line, excludes) != "" {
				continue
			}
			if strings.Index(line, matched) != -1 {
				line = line[strings.Index(line, matched):]
			}
			if strings.Index(line, whitespace) != -1 {
				line = line[:strings.Index(line, whitespace)]
			}
			return line
		}
	}

	return ""
}
