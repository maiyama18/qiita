package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}

func New(logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse("https://qiita.com/api/v2")
	if err != nil {
		return nil, err
	}

	discardLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	return &Client{
		URL:        baseURL,
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}, nil
}

func (c *Client) GetUser(userID string) (*User, error) {
	reqPath := path.Join(c.URL.Path, "users", userID)
	c.URL.Path = reqPath
	c.Logger.Printf("send get request to %s\n", c.URL.String())
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("got response from %s\n", c.URL.String())

	user := User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetPost(postID string) (*Post, error) {
	reqPath := path.Join(c.URL.Path, "items", postID)
	c.URL.Path = reqPath
	c.Logger.Printf("send get request to %s\n", c.URL.String())
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("got response from %s\n", c.URL.String())

    post := Post{}
    if err := json.Unmarshal(body, &post); err != nil {
    	return nil, err
	}

   	return &post, nil
}
