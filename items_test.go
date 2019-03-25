package qiita

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"
)

func TestClient_GetItems(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "items", "GetItems")

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
			inputPage:    3,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items",
			expectedRawQuery:    "page=3&per_page=2",
			expectedPage:        3,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    100,
			expectedTotalCount:  392649,
			expectedItemsLen:    2,
		},
		{
			desc:         "failure-out_of_range",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			itemsResp, err := cli.GetItems(context.Background(), tt.inputPage, tt.inputPerPage)
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

func TestClient_GetItem(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Tokyo")
	mockFilesBaseDir := path.Join("testdata", "responses", "items", "GetItem")

	tests := []struct {
		desc        string
		inputItemID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod          string
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
			desc:        "success",
			inputItemID: "b4ca1773580317e7112e",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:          http.MethodGet,
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
			desc:        "failure-not_exist",
			inputItemID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, "")
			defer teardown()

			item, err := cli.GetItem(context.Background(), tt.inputItemID)
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

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetItemStockers(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "items", "GetItemStockers")

	tests := []struct {
		desc         string
		inputItemID  string
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
			inputItemID:  "b4ca1773580317e7112e",
			inputPage:    3,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/b4ca1773580317e7112e/stockers",
			expectedRawQuery:    "page=3&per_page=2",
			expectedPage:        3,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    100,
			expectedTotalCount:  289,
			expectedUsersLen:    2,
		},
		{
			desc:         "failure-out_of_range",
			inputItemID:  "b4ca1773580317e7112e",
			inputPage:    101,
			inputPerPage: 2,

			mockResponseHeaderFile: "out_of_range-header",
			mockResponseBodyFile:   "out_of_range-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/b4ca1773580317e7112e/stockers",
			expectedRawQuery:    "page=101&per_page=2",
			expectedErrString:   "page parameter should be",
		},
		{
			desc:         "failure-not_found",
			inputItemID:  "nonexistent",
			inputPage:    3,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/nonexistent/stockers",
			expectedRawQuery:    "page=3&per_page=2",
			expectedErrString:   "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			usersResp, err := cli.GetItemStockers(context.Background(), tt.inputItemID, tt.inputPage, tt.inputPerPage)
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

func TestClient_GetItemComments(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "items", "GetItemComments")

	tests := []struct {
		desc        string
		inputItemID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
		expectedCommentsLen int
	}{
		{
			desc:        "success",
			inputItemID: "b4ca1773580317e7112e",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/b4ca1773580317e7112e/comments",
			expectedCommentsLen: 4,
		},
		{
			desc:        "failure-not_found",
			inputItemID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/nonexistent/comments",
			expectedErrString:   "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			comments, err := cli.GetItemComments(context.Background(), tt.inputItemID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedCommentsLen, len(comments))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}
