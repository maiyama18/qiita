package qiita

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	BaseURL    = "https://qiita.com/api/v2"
	PageMin    = 1
	PageMax    = 100
	PerPageMin = 1
	PerPageMax = 100
)

// Client interacts with qiita API
type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	AccessToken string
	UserAgent   string

	Logger *log.Logger
}

// New returns a Client
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
		URL:        baseURL,
		HTTPClient: http.DefaultClient,

		AccessToken: accessToken,
		UserAgent:   "qiita go-client (github.com/muiscript/qiita)",

		Logger: logger,
	}, nil
}
