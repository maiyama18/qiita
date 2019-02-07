package main

import (
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}
