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

func TestClient_GetTag(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "tags", "GetTag")

	tests := []struct {
		desc       string
		inputTagID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod         string
		expectedRequestPath    string
		expectedRawQuery       string
		expectedErrString      string
		expectedID             string
		expectedIconURL        string
		expectedItemsCount     int
		expectedFollowersCount int
	}{
		{
			desc:       "success",
			inputTagID: "react",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:         http.MethodGet,
			expectedRequestPath:    "/tags/react",
			expectedID:             "React",
			expectedIconURL:        "https://s3-ap-northeast-1.amazonaws.com/qiita-tag-image/c4d0439277f132acce23de37f694617b95af5475/medium.jpg?1513495262",
			expectedItemsCount:     2693,
			expectedFollowersCount: 2403,
		},
		{
			desc:       "failure-not_exist",
			inputTagID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/tags/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			tag, err := cli.GetTag(context.Background(), tt.inputTagID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedID, tag.ID)
				assert.Equal(t, tt.expectedIconURL, tag.IconURL)
				assert.Equal(t, tt.expectedItemsCount, tag.ItemsCount)
				assert.Equal(t, tt.expectedFollowersCount, tag.FollowersCount)
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}
