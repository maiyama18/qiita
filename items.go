package qiita

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Item represents an item published on qiita.
type Item struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Body         string `json:"body"`
	RenderedBody string `json:"rendered_body"`
	Private      bool   `json:"private"`
	Coediting    bool   `json:"coediting"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CommentsCount  int `json:"comments_count"`
	LikesCount     int `json:"likes_count"`
	ReactionsCount int `json:"reactions_count"`

	User     *User      `json:"user"`
	ItemTags []*ItemTag `json:"tags"`
}

// ItemTag represents a tag for a qiita item.
type ItemTag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

// ItemsResponse represents a response from qiita API which includes multiple items.
type ItemsResponse struct {
	Items      []*Item
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

func newItemsResponse(items []*Item, header http.Header, page, perPage int) (*ItemsResponse, error) {
	paginationInfo, err := extractPaginationInfo(header, page, perPage)
	if err != nil {
		return nil, err
	}

	return &ItemsResponse{
		Items:      items,
		PerPage:    paginationInfo.PerPage,
		Page:       paginationInfo.Page,
		FirstPage:  paginationInfo.FirstPage,
		LastPage:   paginationInfo.LastPage,
		TotalCount: paginationInfo.TotalCount,
	}, nil
}

// ItemDraft represents an item to be posted for qiita.
type ItemDraft struct {
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	ItemTags []*ItemTag `json:"tags"`
	Private  bool       `json:"private"`
	Tweet    bool       `json:"tweet"`
}

// GetItem fetches the item having provided itemID.
//
// GET /api/v2/items/:item_id
// document: https://qiita.com/api/v2/docs#get-apiv2itemsitem_id
func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("items", itemID), nil, nil)
	if err != nil {
		return nil, err
	}

	var item Item
	code, _, err := c.doRequest(req, &item)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return &item, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetItems fetches all the items posted on qiita.
//
// GET /api/v2/items
// document: http://qiita.com/api/v2/docs#get-apiv2items
func (c *Client) GetItems(ctx context.Context, page, perPage int) (*ItemsResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, "items", query, nil)
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
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetItemComments fetches the comments posted on provided itemID.
//
// GET /api/v2/items/:item_id/comments
// document: http://qiita.com/api/v2/docs#get-apiv2itemsitem_idcomments
func (c *Client) GetItemComments(ctx context.Context, itemID string) ([]*Comment, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("items", itemID, "comments"), nil, nil)
	if err != nil {
		return nil, err
	}

	var comments []*Comment
	code, _, err := c.doRequest(req, &comments)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return comments, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// GetItemStockers fetches the users who stocked the item having provided itemID.
//
// GET /api/v2/items/:item_id/stockers
// document: http://qiita.com/api/v2/docs#get-apiv2itemsitem_idstockers
func (c *Client) GetItemStockers(ctx context.Context, itemID string, page, perPage int) (*UsersResponse, error) {
	if err := validatePaginationLimit(page, perPage); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("items", itemID, "stockers"), query, nil)
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
		return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// CreateItem publishes the item.
// This method requires authentication.
//
// POST /api/v2/items
// document: http://qiita.com/api/v2/docs#post-apiv2items
func (c *Client) CreateItem(ctx context.Context, title, body string, itemTags []*ItemTag, private, tweet bool) (*Item, error) {
	itemDraft := &ItemDraft{Title: title, Body: body, ItemTags: itemTags, Private: private, Tweet: tweet}
	bodyBytes, err := json.Marshal(itemDraft)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "items", nil, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var item Item
	code, _, err := c.doRequest(req, &item)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusCreated:
		return &item, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	case http.StatusForbidden:
		return nil, fmt.Errorf("forbidden. some required field values may be empty or invalid (status = %d)", code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// UpdateItem update the item having provided itemID.
// This method requires authentication.
//
// PATCH /api/v2/items/:item_id
// document: http://qiita.com/api/v2/docs#patch-apiv2itemsitem_id
func (c *Client) UpdateItem(ctx context.Context, itemID string, title, body string, itemTags []*ItemTag, private, tweet bool) (*Item, error) {
	itemDraft := &ItemDraft{Title: title, Body: body, ItemTags: itemTags, Private: private, Tweet: tweet}
	bodyBytes, err := json.Marshal(itemDraft)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPatch, path.Join("items", itemID), nil, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var item Item
	code, _, err := c.doRequest(req, &item)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
		return &item, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	case http.StatusForbidden:
		return nil, fmt.Errorf("forbidden. some required field values may be empty or invalid (status = %d)", code)
	case http.StatusNotFound:
		return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, code)
	default:
		return nil, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// DeleteItem deletes the item.
// This method requires authentication.
//
// DELETE /api/v2/items/:item_id
// document: https://qiita.com/api/v2/docs#delete-apiv2itemsitem_id
func (c *Client) DeleteItem(ctx context.Context, itemID string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path.Join("items", itemID), nil, nil)
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
	case http.StatusForbidden:
		return fmt.Errorf("forbidden. some required field values may be empty or invalid (status = %d)", code)
	case http.StatusNotFound:
		return fmt.Errorf("item with id '%s' not found (status = %d)", itemID, code)
	default:
		return fmt.Errorf("unknown error (status = %d)", code)
	}
}

// CreateItemComment post comments on the item having provided itemID.
// This method requires authentication.
//
// POST /api/v2/items/:item_id/comments
// document: http://qiita.com/api/v2/docs#post-apiv2itemsitem_idcomments
func (c *Client) CreateItemComment(ctx context.Context, itemID string, body string) (*Comment, error) {
	// TODO: implement
	return nil, nil
}

// IsStockedItem returns true if the authenticated user has stocked the item having provided itemID.
// This method requires authentication.
//
// GET /api/v2/items/:item_id/stock
// document: http://qiita.com/api/v2/docs#get-apiv2itemsitem_idstock
func (c *Client) IsStockedItem(ctx context.Context, itemID string) (bool, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path.Join("items", itemID, "stock"), nil, nil)
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
	case http.StatusNotFound:
		return false, nil
	case http.StatusUnauthorized:
		return false, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", code)
	default:
		return false, fmt.Errorf("unknown error (status = %d)", code)
	}
}

// StockItem add the item having provided itemID to the authenticated user's stock list.
// This method requires authentication.
//
// PUT /api/v2/items/:item_id/stock
// document: http://qiita.com/api/v2/docs#put-apiv2itemsitem_idstock
func (c *Client) StockItem(ctx context.Context, itemID string) error {
	// TODO: implement
	return nil
}

// UnstockItem remove the item having provided itemID from the authenticated user's stock list.
// This method requires authentication.
//
// DELETE /api/v2/items/:item_id/stock
// document: http://qiita.com/api/v2/docs#delete-apiv2itemsitem_idstock
func (c *Client) UnstockItem(ctx context.Context, itemID string) error {
	// TODO: implement
	return nil
}
