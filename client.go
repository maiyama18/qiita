package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
