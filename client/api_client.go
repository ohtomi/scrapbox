package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	DefaultHost       = "https://scrapbox.io"
	DefaultExpiration = 60 * 60 // time.Second
	DefaultUserAgent  = "ScrapboxGoClient/0.1.0"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	Token      string
	Expiration time.Duration
}

func NewClient(url *url.URL, token string, expiration int) (*Client, error) {
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
		Token:      token,
		Expiration: time.Duration(expiration) * time.Second,
	}, nil
}

func (c *Client) ExecQuery(ctx context.Context, project string, tags []string, skip, limit int) (*QueryResult, error) {

	var (
		count int
		pages []string

		v interface{}
	)

	host := (*c.URL).Host
	expiration := c.Expiration
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

func (c *Client) GetPage(ctx context.Context, project, page string) (*Page, error) {

	var (
		v interface{}
	)

	host := (*c.URL).Host
	expiration := c.Expiration
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

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {

	baseURL := *c.URL
	u := fmt.Sprintf("%s/%s", baseURL.String(), spath)

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate http request")
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", DefaultUserAgent)
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

func GetURL(host, project, page string) string {
	return fmt.Sprintf("%s/%s/%s", host, project, encodeURIComponent(page))
}

func encodeURIComponent(component string) string {
	regularEscaped := url.QueryEscape(component)
	rParenUnescaped := strings.Replace(regularEscaped, "%28", "(", -1)
	lParenUnescaped := strings.Replace(rParenUnescaped, "%29", ")", -1)
	plusEscaped := strings.Replace(lParenUnescaped, "+", "%20", -1)
	return plusEscaped
}

func buildQueryPath(project string, tags []string, skip, limit int) string {
	params := fmt.Sprintf("skip=%d&sort=updated&limit=%d&q=%s", skip, limit, encodeURIComponent(strings.Join(tags, " ")))
	if len(tags) == 0 {
		return fmt.Sprintf("api/pages/%s?%s", project, params)
	} else {
		return fmt.Sprintf("api/pages/%s/search/query?%s", project, params)
	}
}

func buildPagePath(project, page string) string {
	return fmt.Sprintf("api/pages/%s/%s", project, encodeURIComponent(page))
}
