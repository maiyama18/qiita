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
	TwitterID  string `json:"twitter_screen_name"`
}

type UsersResponse struct {
	Users      []*User
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

type PaginationInfo struct {
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

func constructUsersResponse(users []*User, info *PaginationInfo) *UsersResponse {
	return &UsersResponse{
		Users:      users,
		PerPage:    info.PerPage,
		Page:       info.Page,
		FirstPage:  info.FirstPage,
		LastPage:   info.LastPage,
		TotalCount: info.TotalCount,
	}
}
