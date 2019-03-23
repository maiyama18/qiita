package qiita

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"strings"
	"testing"
)

func TestClient_GetUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUser")

	tests := []struct {
		desc        string
		inputUserID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod         string
		expectedRequestPath    string
		expectedErrString      string
		expectedID             string
		expectedPermanentID    int
		expectedGithubID       string
		expectedPostsCount     int
		expectedFollowersCount int
	}{
		{
			desc:        "success",
			inputUserID: "muiscript",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:         http.MethodGet,
			expectedRequestPath:    "/users/muiscript",
			expectedID:             "muiscript",
			expectedPermanentID:    159260,
			expectedGithubID:       "muiscript",
			expectedPostsCount:     14,
			expectedFollowersCount: 11,
		},
		{
			desc:        "failure-not_exist",
			inputUserID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, "")
			defer teardown()

			user, err := cli.GetUser(context.Background(), tt.inputUserID)
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

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_IsFollowingUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "IsFollowingUser")

	tests := []struct {
		desc        string
		inputUserID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedIsFollowing bool
		expectedErrString   string
	}{
		{
			desc:        "success-following",
			inputUserID: "mizchi",

			mockResponseHeaderFile: "following-header",
			mockResponseBodyFile:   "following-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/mizchi/following",
			expectedIsFollowing: true,
		},
		{
			desc:        "success-not_following",
			inputUserID: "yaotti",

			mockResponseHeaderFile: "not_following-header",
			mockResponseBodyFile:   "not_following-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/yaotti/following",
			expectedIsFollowing: false,
		},
		{
			desc:        "failure-no_token",
			inputUserID: "mizchi",

			mockResponseHeaderFile: "no_token-header",
			mockResponseBodyFile:   "no_token-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/mizchi/following",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, "")
			defer teardown()

			isFollowing, err := cli.IsFollowingUser(context.Background(), tt.inputUserID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedIsFollowing, isFollowing)
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_GetUserFollowees(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUserFollowees")

	tests := []struct {
		desc         string
		inputUserID  string
		inputPage    int
		inputPerPage int

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
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
			inputUserID:  "muiscript",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/followees",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    3,
			expectedTotalCount:  5,
			expectedUsersLen:    2,
		},
		{
			desc:         "failure-page_out_of_range",
			inputUserID:  "muiscript",
			inputPage:    0,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/followees",
			expectedRawQuery:    "page=0&per_page=2",
			expectedErrString:   "page parameter should be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			usersResp, err := cli.GetUserFollowees(context.Background(), tt.inputUserID, tt.inputPage, tt.inputPerPage)
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

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetUserFollowers(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUserFollowers")

	tests := []struct {
		desc         string
		inputUserID  string
		inputPage    int
		inputPerPage int

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
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
			inputUserID:  "muiscript",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
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
			desc:         "failure_page_less_than_100",
			inputUserID:  "muiscript",
			inputPage:    0,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/followers",
			expectedRawQuery:    "page=0&per_page=2",
			expectedErrString:   "page parameter should be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			usersResp, err := cli.GetUserFollowers(context.Background(), tt.inputUserID, tt.inputPage, tt.inputPerPage)
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

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}
