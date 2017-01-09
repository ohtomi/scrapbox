package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
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
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
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

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

type Page struct {
	Title             string
	Lines             []string
	Links             []string
	RelatedPageTitles []string
}

func (c *Client) GetPage(ctx context.Context, page string) (*Page, error) {
	spath := fmt.Sprintf("%s/%s", apiEndpoint, page)
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check status code here…

	var v interface{}
	if err := c.decodeBody(res, &v); err != nil {
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
