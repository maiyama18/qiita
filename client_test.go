package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc   string
		logger *log.Logger

		expectedURL string
		expectedLogger *log.Logger
	}{
		{
			desc: "success",
			logger: log.New(ioutil.Discard, "", log.LstdFlags),

			expectedURL: BASE_URL,
			expectedLogger: log.New(ioutil.Discard, "", log.LstdFlags),
		},
		{
			desc: "success_with_no_logger",
			logger: nil,

			expectedURL: BASE_URL,
			expectedLogger: log.New(ioutil.Discard, "", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, err := New(tt.logger)
			if !assert.Nil(t, err) {
				t.FailNow()
			}

			assert.Equal(t, cli.URL.String(), tt.expectedURL)
			assert.Equal(t, cli.Logger, tt.expectedLogger)
		})
	}
}

func TestClient_GetUser(t *testing.T) {
	tests := []struct {
		desc           string
		id             string
		responseFile   string
		responseStatus int

		expectedRequestPath    string
		expectedErrString      string
		expectedID             string
		expectedPermanentID    int
		expectedGithubID       string
		expectedPostsCount     int
		expectedFollowersCount int
	}{
		{
			desc:           "success",
			id:             "muiscript",
			responseFile:   "users_muiscript",
			responseStatus: http.StatusOK,

			expectedRequestPath:    "/users/muiscript",
			expectedID:             "muiscript",
			expectedPermanentID:    159260,
			expectedGithubID:       "muiscript",
			expectedPostsCount:     14,
			expectedFollowersCount: 11,
		},
		{
			desc:           "failure_nonexistent_user",
			id:             "nonexistent",
			responseFile:   "users_nonexistent",
			responseStatus: http.StatusNotFound,

			expectedRequestPath: "/users/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tt.expectedRequestPath, req.URL.Path)

				dataPath := fmt.Sprintf("./testdata/responses/%s", tt.responseFile)
				f, err := os.Open(dataPath)
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				b, err := ioutil.ReadAll(f)
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				w.WriteHeader(tt.responseStatus)
				w.Write(b)
			}))
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.Nil(t, err)
			cli := &Client{
				URL:        serverURL,
				HTTPClient: server.Client(),
				Logger:     log.New(ioutil.Discard, "", 0),
			}

			user, err := cli.GetUser(context.Background(), tt.id)
			if tt.responseStatus == http.StatusOK {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedID, user.ID)
				assert.Equal(t, tt.expectedPermanentID, user.PermanentID)
				assert.Equal(t, tt.expectedGithubID, user.GithubID)
				assert.Equal(t, tt.expectedGithubID, user.GithubID)
				assert.Equal(t, tt.expectedPostsCount, user.PostsCount)
				assert.Equal(t, tt.expectedFollowersCount, user.FollowersCount)
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_GetItem(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Tokyo")

	tests := []struct {
		desc           string
		id             string
		responseFile   string
		responseStatus int

		expectedRequestPath     string
		expectedErrString       string
		expectedID              string
		expectedTitle           string
		expectedURL             string
		expectedBody            string
		expectedRenderedBody    string
		expectedPrivate         bool
		expectedCreatedAt       time.Time
		expectedUpdatedAt       time.Time
		expectedLikesCount      int
		expectedUserID          string
		expectedUserPermanentID int
		expectedTagNames        []string
	}{
		{
			desc:           "success",
			id:             "b4ca1773580317e7112e",
			responseFile:   "items_b4ca1773580317e7112e",
			responseStatus: http.StatusOK,

			expectedRequestPath:     "/items/b4ca1773580317e7112e",
			expectedID:              "b4ca1773580317e7112e",
			expectedTitle:           "react-router@v4を使ってみよう：シンプルなtutorial",
			expectedURL:             "https://qiita.com/muiscript/items/b4ca1773580317e7112e",
			expectedBody:            "`React`でルーティングをするためのライブラリである`react-router`のv4の基本的な使い方を覚えるために、簡単なwebページを作ってみます。",
			expectedRenderedBody:    "<p><code>React</code>でルーティングをするためのライブラリである<code>react-router</code>のv4の基本的な使い方を覚えるために、簡単なwebページを作ってみます。</p>",
			expectedPrivate:         false,
			expectedCreatedAt:       time.Date(2017, 06, 27, 15, 36, 55, 0, location),
			expectedUpdatedAt:       time.Date(2019, 1, 3, 14, 30, 25, 0, location),
			expectedLikesCount:      309,
			expectedUserID:          "muiscript",
			expectedUserPermanentID: 159260,
		},
		{
			desc:           "failure_nonexistent_item",
			id:             "nonexistent",
			responseFile:   "items_nonexistent",
			responseStatus: http.StatusNotFound,

			expectedRequestPath: "/items/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.responseStatus, tt.expectedRequestPath)
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			if !assert.Nil(t, err) {
				t.FailNow()
			}
			cli := &Client{
				URL:        serverURL,
				HTTPClient: server.Client(),
				Logger:     log.New(ioutil.Discard, "", 0),
			}

			item, err := cli.GetItem(context.Background(), tt.id)
			if tt.responseStatus == http.StatusOK {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedID, item.ID)
				assert.Equal(t, tt.expectedTitle, item.Title)
				assert.Equal(t, tt.expectedURL, item.URL)
				assert.True(t, strings.Contains(item.Body, tt.expectedBody))
				assert.True(t, strings.Contains(item.RenderedBody, tt.expectedRenderedBody))
				assert.Equal(t, tt.expectedPrivate, item.Private)
				assert.True(t, item.CreatedAt.Equal(tt.expectedCreatedAt))
				assert.True(t, item.UpdatedAt.Equal(tt.expectedUpdatedAt))
				assert.Equal(t, tt.expectedLikesCount, item.LikesCount)
				assert.Equal(t, tt.expectedUserID, item.User.ID)
				assert.Equal(t, tt.expectedUserPermanentID, item.User.PermanentID)

			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
			}

		})
	}
}

func newTestServer(t *testing.T, responseFile string, responseStatus int, expectedRequestPath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedRequestPath, req.URL.Path)

		dataPath := fmt.Sprintf("./testdata/responses/%s", responseFile)
		f, err := os.Open(dataPath)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		b, err := ioutil.ReadAll(f)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		w.WriteHeader(responseStatus)
		w.Write(b)
	}))
}
