package qiita

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

// User represents a qiita user.
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

// UsersResponse represents a response from qiita API which includes multiple users.
type UsersResponse struct {
	Users      []*User
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

func newUsersResponse(users []*User, header http.Header, page, perPage int) (*UsersResponse, error) {
	paginationInfo, err := extractPaginationInfo(header, page, perPage)
	if err != nil {
		return nil, err
	}

	return &UsersResponse{
		Users:      users,
		PerPage:    paginationInfo.PerPage,
		Page:       paginationInfo.Page,
		FirstPage:  paginationInfo.FirstPage,
		LastPage:   paginationInfo.LastPage,
		TotalCount: paginationInfo.TotalCount,
	}, nil
}

type paginationInfo struct {
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

// GetUser fetches the user having provided userID.
//
// GET /api/v2/users/:user_id
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_id
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID), nil, nil)
	if err != nil {
		return nil, err
	}

	var user User
	code, _, err := c.doRequest(req, &user)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return &user, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUsers fetches all the users.
// The number of users included in one response and the page number should be provided.
//
// GET /api/v2/users
// document: https://qiita.com/api/v2/docs#get-apiv2users
func (c *Client) GetUsers(ctx context.Context, page, perPage int) (*UsersResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, "users", query, nil)
	if err != nil {
		return nil, err
	}

	var users []*User
	code, header, err := c.doRequest(req, &users)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newUsersResponse(users, header, page, perPage)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUserFollowees fetches all the followees of the user having provided userID.
// The number of users included in one response and the page number should be provided.
//
// GET /api/v2/users/:user_id/followees
// document: http://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowees
func (c *Client) GetUserFollowees(ctx context.Context, userID string, page, perPage int) (*UsersResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "followees"), query, nil)
	if err != nil {
		return nil, err
	}

	var users []*User
	code, header, err := c.doRequest(req, &users)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newUsersResponse(users, header, page, perPage)
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUserFollowers fetches all the followers of the user having provided userID.
// The number of users included in one response and the page number should be provided.
//
// GET /api/v2/users/:user_id/followers
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowers
func (c *Client) GetUserFollowers(ctx context.Context, userID string, page, perPage int) (*UsersResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "followers"), query, nil)
	if err != nil {
		return nil, err
	}

	var users []*User
	code, header, err := c.doRequest(req, &users)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newUsersResponse(users, header, page, perPage)
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUserItems fetches the items created by the user having provided userID.
//
// GET /api/v2/users/:user_id/items
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_iditems
func (c *Client) GetUserItems(ctx context.Context, userID string, page, perPage int) (*ItemsResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "items"), query, nil)
	if err != nil {
		return nil, err
	}

	var items []*Item
	code, header, err := c.doRequest(req, &items)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newItemsResponse(items, header, page, perPage)
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUserStocks fetches the items stocked by the user having provided userID.
//
// GET /api/v2/users/:user_id/stocks
// document: http://qiita.com/api/v2/docs#get-apiv2usersuser_idstocks
func (c *Client) GetUserStocks(ctx context.Context, userID string, page, perPage int) (*ItemsResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "stocks"), query, nil)
	if err != nil {
		return nil, err
	}

	var items []*Item
	code, header, err := c.doRequest(req, &items)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newItemsResponse(items, header, page, perPage)
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetUserFollowingTags fetches the tags followed by the user having provided userID.
//
// GET /api/v2/users/:user_id/following_tags
// document: http://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowing_tags
func (c *Client) GetUserFollowingTags(ctx context.Context, userID string, page, perPage int) (*TagsResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "following_tags"), query, nil)
	if err != nil {
		return nil, err
	}

	var tags []*Tag
	code, header, err := c.doRequest(req, &tags)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newTagsResponse(tags, header, page, perPage)
	case http.StatusNotFound:
		return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// IsFollowingUser returns true if the authenticated user is following the user having provided userID.
// This method requires authentication.
//
// GET /api/v2/users/:user_id/following
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowing
func (c *Client) IsFollowingUser(ctx context.Context, userID string) (bool, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("users", userID, "following"), nil, nil)
	if err != nil {
		return false, err
	}

	code, _, err := c.doRequest(req, &struct{}{})
	if err != nil {
		return false, err
	}
	switch code {
	case http.StatusNoContent:
		return true, nil
	case http.StatusUnauthorized:
		return false, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// FollowUser follows the user having provided userID.
// This method requires authentication.
//
// PUT /api/v2/users/:user_id/following
// document: http://qiita.com/api/v2/docs#put-apiv2usersuser_idfollowing
func (c *Client) FollowUser(ctx context.Context, userID string) error {
	req, err := c.newRequest(ctx, http.MethodPut, path.Join("users", userID, "following"), nil, nil)
	if err != nil {
		return err
	}

	code, _, err := c.doRequest(req, &struct{}{})
	if err != nil {
		return err
	}
	switch code {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	case http.StatusNotFound:
		return fmt.Errorf("not found. user with id '%s' does not exist (status = %d)", userID, code)
	case http.StatusForbidden:
		return fmt.Errorf("forbidden. you may already have followed user with id '%s' (status = %d)", userID, code)
	default:
		return fmt.Errorf("unknown error (status = %d)", code)
	}
}

// UnfollowUser unfollows the user having provided userID.
// This method requires authentication.
//
// DELETE /api/v2/users/:user_id/following
// document: http://qiita.com/api/v2/docs#delete-apiv2usersuser_idfollowing
func (c *Client) UnfollowUser(ctx context.Context, userID string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path.Join("users", userID, "following"), nil, nil)
	if err != nil {
		return err
	}

	code, _, err := c.doRequest(req, &struct{}{})
	if err != nil {
		return err
	}
	switch code {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	case http.StatusNotFound:
		return fmt.Errorf("not found. user with id '%s' does not exist (status = %d)", userID, code)
	case http.StatusForbidden:
		return fmt.Errorf("forbidden. you may already have not followed user with id '%s' (status = %d)", userID, code)
	default:
		return fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetAuthenticatedUser returns the user who is associated with provided access token.
// This method requires authentication.
//
// GET /api/v2/authenticated_user
// document: http://qiita.com/api/v2/docs#get-apiv2authenticated_user
func (c *Client) GetAuthenticatedUser(ctx context.Context) (*User, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("authenticated_user"), nil, nil)
	if err != nil {
		return nil, err
	}

	var user User
	code, _, err := c.doRequest(req, &user)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return &user, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetAuthenticatedUserItems fetches the item created by the authenticated user.
// This method requires authentication.
//
// GET /api/v2/authenticated_user/items
// document: http://qiita.com/api/v2/docs#get-apiv2authenticated_useritems
func (c *Client) GetAuthenticatedUserItems(ctx context.Context, page, perPage int) (*ItemsResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("authenticated_user", "items"), query, nil)
	if err != nil {
		return nil, err
	}

	var items []*Item
	code, header, err := c.doRequest(req, &items)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return newItemsResponse(items, header, page, perPage)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}
