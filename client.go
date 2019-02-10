package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

const BASE_URL = "https://qiita.com/api/v2"

type Client struct {
	URL         *url.URL
	HTTPClient  *http.Client
	AccessToken string
	Logger      *log.Logger
}

func New(accessToken string, logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse(BASE_URL)
	if err != nil {
		return nil, err
	}

	discardLogger := log.New(ioutil.Discard, "", 0)
	if logger == nil {
		logger = discardLogger
	}

	return &Client{
		URL:         baseURL,
		HTTPClient:  http.DefaultClient,
		AccessToken: accessToken,
		Logger:      logger,
	}, nil
}

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

func (c *Client) GetUsers(ctx context.Context, page int, perPage int) (*UsersResponse, error) {
	if page < 1 || 100 < page {
		return nil, fmt.Errorf("page parameter should be between 1 and 100. got %d", page)
	}
	if perPage < 1 || 100 < perPage {
		return nil, fmt.Errorf("perPage parameter should be between 1 and 100. got %d", perPage)
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", "users", query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

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

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID), nil, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, resp.StatusCode)
		default:
			return nil, fmt.Errorf("unknown error (status = %d)", resp.StatusCode)
		}
	}

	var user User
	if err := c.decodeBody(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetFollowees(ctx context.Context, userID string, page int, perPage int) (*UsersResponse, error) {
	if page < 1 || 100 < page {
		return nil, fmt.Errorf("page parameter should be between 1 and 100. got %d", page)
	}
	if perPage < 1 || 100 < perPage {
		return nil, fmt.Errorf("perPage parameter should be between 1 and 100. got %d", perPage)
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "followees"), query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

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

func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("items", itemID), nil, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, resp.StatusCode)
		default:
			return nil, fmt.Errorf("unknown error (status = %d)", resp.StatusCode)
		}
	}

	var item Item
	if err := c.decodeBody(resp, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) IsFollowingUser(ctx context.Context, userID string) (bool, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "following"), nil, nil)
	if err != nil {
		return false, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else {
		return false, nil
	}
}
