package main

type User struct {
	ID           string `json:"id"`
	PermanentID  int    `json:"permanent_id"`
	ImageURL     string `json:"profile_image_url"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	WebsiteURL   string `json:"website_url"`
	Organization string `json:"organization"`
	TeamOnly     bool   `json:"team_only"`

	PostsCount     int `json:"items_count"`
	FolloweesCount int `json:"followees_count"`
	FollowersCount int `json:"followers_count"`

	GithubID   string `json:"github_login_name"`
	LinkedinID string `json:"linkedin_id"`
	twitterID  string `json: "twitter_screen_name"`
}
