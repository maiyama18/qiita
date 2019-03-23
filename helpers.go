package qiita

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func (c *Client) newRequest(ctx context.Context, method string, relativePath string, query map[string]string, body io.Reader) (*http.Request, error) {
	reqUrl := *c.URL
	reqUrl.Path = path.Join(reqUrl.Path, relativePath)

	if query != nil {
		q := reqUrl.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		reqUrl.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "qiita go-client (github.com/muiscript/qiita)")
	if c.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	return req, nil
}

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) parseHeaderLink(resp *http.Response) (map[string]*url.URL, error) {
	links := make(map[string]*url.URL)

	linksStr := resp.Header.Get("link")
	rx := regexp.MustCompile("<(.*)>.*rel=\"(.*)\"")

	for _, link := range strings.Split(linksStr, ", ") {
		m := rx.FindStringSubmatch(link)

		rel := m[2]
		linkURL, err := url.Parse(m[1])
		if err != nil {
			return nil, err
		}

		links[rel] = linkURL
	}

	return links, nil
}

func (c *Client) validatePaginationLimit(page, pageMin, pageMax, perPage, perPageMin, perPageMax int) error {
	if page < pageMin || pageMax < page {
		return fmt.Errorf("page parameter should be between 1 and 100. got %d", page)
	}
	if perPage < perPageMin || perPageMax < perPage {
		return fmt.Errorf("perPage parameter should be between 1 and 100. got %d", perPage)
	}
	return nil
}

func (c *Client) extractPaginationInfo(resp *http.Response, page int, perPage int) (*paginationInfo, error) {
	links, err := c.parseHeaderLink(resp)
	if err != nil {
		return nil, err
	}
	lastURL := links["last"]
	lastPage, err := strconv.Atoi(lastURL.Query().Get("page"))
	if err != nil {
		return nil, err
	}
	if lastPage > 100 {
		lastPage = 100
	}

	totalCount, err := strconv.Atoi(resp.Header.Get("total-count"))
	if err != nil {
		return nil, err
	}

	return &paginationInfo{
		Page:       page,
		PerPage:    perPage,
		FirstPage:  1,
		LastPage:   lastPage,
		TotalCount: totalCount,
	}, nil
}
