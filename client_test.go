package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc        string
		accessToken string
		logger      *log.Logger

		expectedURL    string
		expectedLogger *log.Logger
	}{
		{
			desc:        "success",
			accessToken: "access_token",
			logger:      log.New(os.Stdout, "", log.LstdFlags),

			expectedURL:    BASE_URL,
			expectedLogger: log.New(os.Stdout, "", log.LstdFlags),
		},
		{
			desc:        "success_with_no_logger",
			accessToken: "access_token",
			logger:      nil,

			expectedURL:    BASE_URL,
			expectedLogger: log.New(ioutil.Discard, "", 0),
		},
		{
			desc:        "success_with_no_access_token",
			accessToken: "",
			logger:      log.New(os.Stdout, "", log.LstdFlags),

			expectedURL:    BASE_URL,
			expectedLogger: log.New(os.Stdout, "", log.LstdFlags),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, err := New(tt.accessToken, tt.logger)
			if !assert.Nil(t, err) {
				t.FailNow()
			}

			assert.Equal(t, cli.URL.String(), tt.expectedURL)
			assert.Equal(t, cli.Logger, tt.expectedLogger)
		})
	}
}

func TestClient_GetFollowees(t *testing.T) {
	tests := []struct {
		desc         string
		userID       string
		page         int
		perPage      int
		responseFile string

		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
		expectedPage        int
		expectedPerPage     int
		expectedFirstPage   int
		expectedLastPage    int
		expectedTotalCount  int
		expectedUsersLen    int
	}{
		{
			desc:         "success",
			userID:       "muiscript",
			page:         2,
			perPage:      2,
			responseFile: "users_muiscript_followees?page=2&per_page=2",

			expectedRequestPath: "/users/muiscript/followees",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    6,
			expectedTotalCount:  11,
			expectedUsersLen:    2,
		},
		{
			desc:         "success_page_larger_than_last",
			userID:       "muiscript",
			page:         10,
			perPage:      2,
			responseFile: "users_muiscript_followees?page=10&per_page=2",

			expectedRequestPath: "/users/muiscript/followees",
			expectedRawQuery:    "page=10&per_page=2",
			expectedPage:        10,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    6,
			expectedTotalCount:  11,
			expectedUsersLen:    0,
		},
		{
			desc:         "failure_page_less_than_100",
			userID:       "muiscript",
			page:         0,
			perPage:      2,
			responseFile: "users_muiscript_followees?page=0&per_page=2",

			expectedRequestPath: "/users/muiscript/followees",
			expectedRawQuery:    "page=0&per_page=2",
			expectedErrString:   "page parameter should be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedRequestPath, tt.expectedRawQuery)
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.Nil(t, err)
			cli := &Client{
				URL:        serverURL,
				HTTPClient: server.Client(),
				Logger:     log.New(ioutil.Discard, "", 0),
			}

			usersResp, err := cli.GetFollowees(context.Background(), tt.userID, tt.page, tt.perPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, usersResp.Page)
				assert.Equal(t, tt.expectedPerPage, usersResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, usersResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, usersResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, usersResp.TotalCount)
				assert.Equal(t, tt.expectedUsersLen, len(usersResp.Users))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetFollowers(t *testing.T) {
	tests := []struct {
		desc         string
		userID       string
		page         int
		perPage      int
		responseFile string

		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
		expectedPage        int
		expectedPerPage     int
		expectedFirstPage   int
		expectedLastPage    int
		expectedTotalCount  int
		expectedUsersLen    int
	}{
		{
			desc:         "success",
			userID:       "muiscript",
			page:         2,
			perPage:      2,
			responseFile: "users_muiscript_followers?page=2&per_page=2",

			expectedRequestPath: "/users/muiscript/followers",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    6,
			expectedTotalCount:  11,
			expectedUsersLen:    2,
		},
		{
			desc:         "success_page_larger_than_last",
			userID:       "muiscript",
			page:         10,
			perPage:      2,
			responseFile: "users_muiscript_followers?page=10&per_page=2",

			expectedRequestPath: "/users/muiscript/followers",
			expectedRawQuery:    "page=10&per_page=2",
			expectedPage:        10,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    6,
			expectedTotalCount:  11,
			expectedUsersLen:    0,
		},
		{
			desc:         "failure_page_less_than_100",
			userID:       "muiscript",
			page:         0,
			perPage:      2,
			responseFile: "users_muiscript_followers?page=0&per_page=2",

			expectedRequestPath: "/users/muiscript/followers",
			expectedRawQuery:    "page=0&per_page=2",
			expectedErrString:   "page parameter should be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedRequestPath, tt.expectedRawQuery)
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.Nil(t, err)
			cli := &Client{
				URL:        serverURL,
				HTTPClient: server.Client(),
				Logger:     log.New(ioutil.Discard, "", 0),
			}

			usersResp, err := cli.GetFollowers(context.Background(), tt.userID, tt.page, tt.perPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, usersResp.Page)
				assert.Equal(t, tt.expectedPerPage, usersResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, usersResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, usersResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, usersResp.TotalCount)
				assert.Equal(t, tt.expectedUsersLen, len(usersResp.Users))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetUser(t *testing.T) {
	tests := []struct {
		desc         string
		id           string
		responseFile string

		expectedRequestPath    string
		expectedErrString      string
		expectedID             string
		expectedPermanentID    int
		expectedGithubID       string
		expectedPostsCount     int
		expectedFollowersCount int
	}{
		{
			desc:         "success",
			id:           "muiscript",
			responseFile: "users_muiscript",

			expectedRequestPath:    "/users/muiscript",
			expectedID:             "muiscript",
			expectedPermanentID:    159260,
			expectedGithubID:       "muiscript",
			expectedPostsCount:     14,
			expectedFollowersCount: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedRequestPath, "")
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.Nil(t, err)
			cli := &Client{
				URL:        serverURL,
				HTTPClient: server.Client(),
				Logger:     log.New(ioutil.Discard, "", 0),
			}

			user, err := cli.GetUser(context.Background(), tt.id)
			if tt.expectedErrString == "" {
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

func TestClient_IsFollowingUser(t *testing.T) {
	tests := []struct {
		desc         string
		targetUserID string
		responseFile string

		expectedRequestPath string
		expectedIsFollowing bool
		expectedErrString   string
	}{
		{
			desc:         "success_following",
			targetUserID: "mizchi",
			responseFile: "users_mizchi_following",

			expectedRequestPath: "/users/mizchi/following",
			expectedIsFollowing: true,
		},
		{
			desc:         "success_not_following",
			targetUserID: "yaotti",
			responseFile: "users_yaotti_following",

			expectedRequestPath: "/users/yaotti/following",
			expectedIsFollowing: false,
		},
		{
			desc:         "failure_no_token",
			targetUserID: "mizchi",
			responseFile: "users_mizchi_following-no_token",

			expectedRequestPath: "/users/mizchi/following",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedRequestPath, "")
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

			isFollowing, err := cli.IsFollowingUser(context.Background(), tt.targetUserID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedIsFollowing, isFollowing)
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
		desc         string
		id           string
		responseFile string

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
			desc:         "success",
			id:           "b4ca1773580317e7112e",
			responseFile: "items_b4ca1773580317e7112e",

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
			desc:         "failure_nonexistent_item",
			id:           "nonexistent",
			responseFile: "items_nonexistent",

			expectedRequestPath: "/items/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedRequestPath, "")
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
			if tt.expectedErrString == "" {
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

func newTestServer(t *testing.T, responseFile string, expectedRequestPath string, expectedRawQuery string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !assert.Equal(t, expectedRequestPath, req.URL.Path) {
			t.FailNow()
		}
		if !assert.Equal(t, expectedRawQuery, req.URL.RawQuery) {
			t.FailNow()
		}

		headerPath := fmt.Sprintf("./testdata/responses/%s-header", responseFile)
		h, err := os.Open(headerPath)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		sc := bufio.NewScanner(h)

		var statusCode int
		for sc.Scan() {
			line := sc.Text()
			if len(line) == 0 {
				continue
			}

			if strings.HasPrefix(line, "HTTP/2") {
				codeStr := strings.Split(line, " ")[1]
				statusCode, _ = strconv.Atoi(codeStr)
			} else {
				key := strings.Split(line, ": ")[0]
				value := strings.Split(line, ": ")[1]
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(statusCode)

		bodyPath := fmt.Sprintf("./testdata/responses/%s-body", responseFile)
		f, err := os.Open(bodyPath)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		b, err := ioutil.ReadAll(f)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		_, err = w.Write(b)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}))
}
