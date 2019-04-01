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

func (c *Client) newRequest(ctx context.Context, method string, relativePath string, queries map[string]string, headers map[string]string, body io.Reader) (*http.Request, error) {
	reqUrl := *c.URL
	reqUrl.Path = path.Join(reqUrl.Path, relativePath)

	if queries != nil {
		q := reqUrl.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		reqUrl.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["User-Agent"] = c.UserAgent
	if c.AccessToken != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, body interface{}) (int, http.Header, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return resp.StatusCode, resp.Header, nil
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

func validatePaginationLimit(page, perPage int) error {
	if page < PageMin || PageMax < page {
		return fmt.Errorf("page parameter should be between %d and %d. got %d", PageMin, PageMax, page)
	}
	if perPage < PerPageMin || PerPageMax < perPage {
		return fmt.Errorf("perPage parameter should be between %d and %d. got %d", PerPageMin, PerPageMax, perPage)
	}
	return nil
}

func extractPaginationInfo(header http.Header, page int, perPage int) (*paginationInfo, error) {
	links, err := parseHeaderLink(header)
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

func parseHeaderLink(header http.Header) (map[string]*url.URL, error) {
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
