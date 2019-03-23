package main

import (
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const BaseURL = "https://qiita.com/api/v2"

type Client struct {
	URL         *url.URL
	HTTPClient  *http.Client
	AccessToken string
	Logger      *log.Logger
}

func New(accessToken string, logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse(BaseURL)
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

func (c *Client) GetUsers(ctx context.Context, page int, perPage int) (*UsersResponse, error) {
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
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

	return c.extractUsersResponse(resp, page, perPage)
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
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
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

	return c.extractUsersResponse(resp, page, perPage)
}

func (c *Client) GetFollowers(ctx context.Context, userID string, page int, perPage int) (*UsersResponse, error) {
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "followers"), query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return c.extractUsersResponse(resp, page, perPage)
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
