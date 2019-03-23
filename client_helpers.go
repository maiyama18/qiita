package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
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
		url, err := url.Parse(m[1])
		if err != nil {
			return nil, err
		}

		links[rel] = url
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

func (c *Client) extractUsersResponse(resp *http.Response, page int, perPage int) (*UsersResponse, error) {
	var users []*User
	if err := c.decodeBody(resp, &users); err != nil {
		return nil, err
	}

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

	usersResp := &UsersResponse{
		Page:       page,
		PerPage:    perPage,
		FirstPage:  1,
		LastPage:   lastPage,
		TotalCount: totalCount,
		Users:      users,
	}

	return usersResp, nil
}

func (c *Client) extractPaginationInfo(resp *http.Response, page int, perPage int) (*PaginationInfo, error) {
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

	return &PaginationInfo{
		Page:       page,
		PerPage:    perPage,
		FirstPage:  1,
		LastPage:   lastPage,
		TotalCount: totalCount,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method string, relativePath string, query map[string]string, body io.Reader) (*http.Request, error) {
	url := *c.URL
	url.Path = path.Join(url.Path, relativePath)

	if query != nil {
		q := url.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, url.String(), body)
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
