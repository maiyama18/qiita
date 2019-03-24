package qiita

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	req.Header.Set("User-Agent", c.UserAgent)
	if c.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, body interface{}) (int, http.Header, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	c.Logger.Printf("send %s request to %s\n", req.Method, c.URL.String())

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, body); err != nil {
			return 0, nil, err
		}
	}

	return resp.StatusCode, resp.Header, nil
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

func (c *Client) extractPaginationInfo(header http.Header, page int, perPage int) (*paginationInfo, error) {
	links, err := c.parseHeaderLink(header)
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

	totalCount, err := strconv.Atoi(header.Get("total-count"))
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

func (c *Client) parseHeaderLink(header http.Header) (map[string]*url.URL, error) {
	links := make(map[string]*url.URL)

	linksStr := header.Get("link")
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
