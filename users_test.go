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

func TestClient_GetUsers(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUsers")

	tests := []struct {
		desc         string
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
			inputPage:    3,
			inputPerPage: 20,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users",
			expectedRawQuery:    "page=3&per_page=20",
			expectedPage:        3,
			expectedPerPage:     20,
			expectedFirstPage:   1,
			expectedLastPage:    100,
			expectedTotalCount:  326706,
			expectedUsersLen:    20,
		},
		{
			desc:         "failure-out_of_range",
			inputPage:    3,
			inputPerPage: 101,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users",
			expectedRawQuery:    "page=3&per_page=101",
			expectedErrString:   "perPage parameter should be",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			usersResp, err := cli.GetUsers(context.Background(), tt.inputPage, tt.inputPerPage)
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

func TestClient_GetUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUser")

	tests := []struct {
		desc        string
		inputUserID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod         string
		expectedRequestPath    string
		expectedRawQuery       string
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
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
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
		expectedRawQuery    string
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
			desc:        "failure-not_exist",
			inputUserID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/following",
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
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
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
		{
			desc:         "failure-not_exist",
			inputUserID:  "nonexistent",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/followees",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "not found",
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
			desc:         "failure-page_out_of_range",
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
		{
			desc:         "failure-not_exist",
			inputUserID:  "nonexistent",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/followers",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "not found",
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

func TestClient_GetUserItems(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUserItems")

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
		expectedItemsLen    int
	}{
		{
			desc:         "success",
			inputUserID:  "muiscript",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/items",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    7,
			expectedTotalCount:  14,
			expectedItemsLen:    2,
		},
		{
			desc:         "failure-page_out_of_range",
			inputUserID:  "muiscript",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/items",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
		{
			desc:         "failure-not_exist",
			inputUserID:  "nonexistent",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/items",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			itemsResp, err := cli.GetUserItems(context.Background(), tt.inputUserID, tt.inputPage, tt.inputPerPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, itemsResp.Page)
				assert.Equal(t, tt.expectedPerPage, itemsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, itemsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, itemsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, itemsResp.TotalCount)
				assert.Equal(t, tt.expectedItemsLen, len(itemsResp.Items))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetUserStocks(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUserStocks")

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
		expectedItemsLen    int
	}{
		{
			desc:         "success",
			inputUserID:  "muiscript",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/stocks",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    11,
			expectedTotalCount:  22,
			expectedItemsLen:    2,
		},
		{
			desc:         "failure-page_out_of_range",
			inputUserID:  "muiscript",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/stocks",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
		{
			desc:         "failure-not_exist",
			inputUserID:  "nonexistent",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/stocks",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			itemsResp, err := cli.GetUserStocks(context.Background(), tt.inputUserID, tt.inputPage, tt.inputPerPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, itemsResp.Page)
				assert.Equal(t, tt.expectedPerPage, itemsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, itemsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, itemsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, itemsResp.TotalCount)
				assert.Equal(t, tt.expectedItemsLen, len(itemsResp.Items))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetUserFollowingTags(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetUserFollowingTags")

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
		expectedTagsLen     int
	}{
		{
			desc:         "success",
			inputUserID:  "muiscript",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/following_tags",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    3,
			expectedTotalCount:  6,
			expectedTagsLen:     2,
		},
		{
			desc:         "failure-page_out_of_range",
			inputUserID:  "muiscript",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/muiscript/following_tags",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
		{
			desc:         "failure-not_exist",
			inputUserID:  "nonexistent",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/users/nonexistent/following_tags",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			tagsResp, err := cli.GetUserFollowingTags(context.Background(), tt.inputUserID, tt.inputPage, tt.inputPerPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, tagsResp.Page)
				assert.Equal(t, tt.expectedPerPage, tagsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, tagsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, tagsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, tagsResp.TotalCount)
				assert.Equal(t, tt.expectedTagsLen, len(tagsResp.Tags))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_FollowUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "FollowUser")

	tests := []struct {
		desc        string
		inputUserID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
	}{
		{
			desc:        "success",
			inputUserID: "yaotti",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodPut,
			expectedRequestPath: "/users/yaotti/following",
		},
		{
			desc:        "failure-already_following",
			inputUserID: "mizchi",

			mockResponseHeaderFile: "already_following-header",
			mockResponseBodyFile:   "already_following-body",

			expectedMethod:      http.MethodPut,
			expectedRequestPath: "/users/mizchi/following",
			expectedErrString:   "forbidden. you may already have followed",
		},
		{
			desc:        "failure-not_exist",
			inputUserID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodPut,
			expectedRequestPath: "/users/nonexistent/following",
			expectedErrString:   "not found",
		},
		{
			desc:        "failure-no_token",
			inputUserID: "yaotti",

			mockResponseHeaderFile: "no_token-header",
			mockResponseBodyFile:   "no_token-body",

			expectedMethod:      http.MethodPut,
			expectedRequestPath: "/users/yaotti/following",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			err := cli.FollowUser(context.Background(), tt.inputUserID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_UnfollowUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "UnfollowUser")

	tests := []struct {
		desc        string
		inputUserID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
	}{
		{
			desc:        "success",
			inputUserID: "mizchi",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodDelete,
			expectedRequestPath: "/users/mizchi/following",
		},
		{
			desc:        "failure-not_following",
			inputUserID: "yaotti",

			mockResponseHeaderFile: "not_following-header",
			mockResponseBodyFile:   "not_following-body",

			expectedMethod:      http.MethodDelete,
			expectedRequestPath: "/users/yaotti/following",
			expectedErrString:   "forbidden. you may already have not followed",
		},
		{
			desc:        "failure-not_exist",
			inputUserID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodDelete,
			expectedRequestPath: "/users/nonexistent/following",
			expectedErrString:   "not found",
		},
		{
			desc:        "failure-no_token",
			inputUserID: "mizchi",

			mockResponseHeaderFile: "no_token-header",
			mockResponseBodyFile:   "no_token-body",

			expectedMethod:      http.MethodDelete,
			expectedRequestPath: "/users/mizchi/following",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			err := cli.UnfollowUser(context.Background(), tt.inputUserID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_GetAuthenticatedUser(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetAuthenticatedUser")

	tests := []struct {
		desc string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod         string
		expectedRequestPath    string
		expectedRawQuery       string
		expectedErrString      string
		expectedID             string
		expectedPermanentID    int
		expectedGithubID       string
		expectedPostsCount     int
		expectedFollowersCount int
	}{
		{
			desc: "success",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:         http.MethodGet,
			expectedRequestPath:    "/authenticated_user",
			expectedID:             "muiscript",
			expectedPermanentID:    159260,
			expectedGithubID:       "muiscript",
			expectedPostsCount:     14,
			expectedFollowersCount: 12,
		},
		{
			desc: "failure-no_token",

			mockResponseHeaderFile: "no_token-header",
			mockResponseBodyFile:   "no_token-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/authenticated_user",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			user, err := cli.GetAuthenticatedUser(context.Background())
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

func TestClient_GetAuthenticatedUserItems(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "users", "GetAuthenticatedUserItems")

	tests := []struct {
		desc         string
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
		expectedItemsLen    int
	}{
		{
			desc:         "success",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/authenticated_user/items",
			expectedRawQuery:    "page=2&per_page=2",
			expectedPage:        2,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    8,
			expectedTotalCount:  16,
			expectedItemsLen:    2,
		},
		{
			desc:         "failure-page_out_of_range",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/authenticated_user/items",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
		{
			desc:         "failure-no_token",
			inputPage:    2,
			inputPerPage: 2,

			mockResponseHeaderFile: "no_token-header",
			mockResponseBodyFile:   "no_token-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/authenticated_user/items",
			expectedRawQuery:    "page=2&per_page=2",
			expectedErrString:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			itemsResp, err := cli.GetAuthenticatedUserItems(context.Background(), tt.inputPage, tt.inputPerPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, itemsResp.Page)
				assert.Equal(t, tt.expectedPerPage, itemsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, itemsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, itemsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, itemsResp.TotalCount)
				assert.Equal(t, tt.expectedItemsLen, len(itemsResp.Items))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}
