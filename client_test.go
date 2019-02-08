package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: curl -i --raw で得たデータを使ってテストを追加する
// 参考： http://hassansin.github.io/Unit-Testing-http-client-in-Go
// type User struct {
// 	ID           string `json:"id"`
// 	PermanentID  int    `json:"permanent_id"`
// 	ImageURL     string `json:"profile_image_url"`
// 	Name         string `json:"name"`
// 	Location     string `json:"location"`
// 	Description  string `json:"description"`
// 	WebsiteURL   string `json:"website_url"`
// 	Organization string `json:"organization"`
// 	TeamOnly     bool   `json:"team_only"`
//
// 	PostsCount     int `json:"items_count"`
// 	FolloweesCount int `json:"followees_count"`
// 	FollowersCount int `json:"followers_count"`
//
// 	GithubID   string `json:"github_login_name"`
// 	LinkedinID string `json:"linkedin_id"`
// 	twitterID  string `json: "twitter_screen_name"`
// }

func TestClient_GetUser(t *testing.T) {
	tests := []struct {
		responseFile        string

		expectedID          string
		expectedPermanentID int
		expectedGithubID    string
		postsCount          int
		FollowersCount      int
	}{
		{
			responseFile:        "user_response",

			expectedID:          "muiscript",
			expectedPermanentID: 159260,
			expectedGithubID:    "muiscript",
			postsCount:          14,
			FollowersCount:      11,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	}))

	c := New()
}
