package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestClient_GetUser(t *testing.T) {
	tests := []struct {
		desc         string
		id           string
		responseFile string

		expectedRequestPath    string
		expectedID             string
		expectedPermanentID    int
		expectedGithubID       string
		expectedPostsCount     int
		expectedFollowersCount int
	}{
		{
			desc:         "success",
			id:           "muiscript",
			responseFile: "user_response",

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
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tt.expectedRequestPath, req.URL.Path)

				f, err := os.Open("./testdata/user_response")
				assert.Nil(t, err)

				b, err := ioutil.ReadAll(f)
				assert.Nil(t, err)

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
			assert.Nil(t, err, "err not nil")

			assert.Equal(t, tt.expectedID, user.ID)
			assert.Equal(t, tt.expectedPermanentID, user.PermanentID)
			assert.Equal(t, tt.expectedGithubID, user.GithubID)
			assert.Equal(t, tt.expectedGithubID, user.GithubID)
			assert.Equal(t, tt.expectedPostsCount, user.PostsCount)
			assert.Equal(t, tt.expectedFollowersCount, user.FollowersCount)
		})
	}
}
