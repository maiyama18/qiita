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
)

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
			desc:           "failure: nonexistent user",
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
