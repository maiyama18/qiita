package qiita

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestClient_GetItem(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Tokyo")

	tests := []struct {
		desc         string
		id           string
		responseFile string

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
			desc:         "success",
			id:           "b4ca1773580317e7112e",
			responseFile: "items_b4ca1773580317e7112e",

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
			desc:         "failure_nonexistent_item",
			id:           "nonexistent",
			responseFile: "items_nonexistent",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/items/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			server := newTestServer(t, tt.responseFile, tt.expectedMethod, tt.expectedRequestPath, "")
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
