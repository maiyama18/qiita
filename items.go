package qiita

import (
	"context"
	"fmt"
	"net/http"
	"path"
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

	User *User      `json:"user"`
	Tags []*ItemTag `json:"tags"`
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
	Title   string `json:"title"`
	Body    string `json:"body"`
	Private bool   `json:"private"`
	Tweet   bool   `json:"tweet"`
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
func (c *Client) GetItems(ctx context.Context) ([]*Item, error) {
	return nil, nil
}

// GetItemComments fetches the comments posted on provided itemID.
//
// GET /api/v2/items/:item_id/comments
// document: http://qiita.com/api/v2/docs#get-apiv2itemsitem_idcomments
func (c *Client) GetItemComments(ctx context.Context, itemID string) ([]*Comment, error) {
	// TODO: implement
	return nil, nil
}

// GetItemStockers fetches the users who stocked the item having provided itemID.
//
// GET /api/v2/items/:item_id/stockers
// document: http://qiita.com/api/v2/docs#get-apiv2itemsitem_idstockers
func (c *Client) GetItemStockers(ctx context.Context, itemID string) ([]*User, error) {
	// TODO: implement
	return nil, nil
}

// CreateItem posts the item.
// This method requires authentication.
//
// POST /api/v2/items
// document: http://qiita.com/api/v2/docs#post-apiv2items
func (c *Client) CreateItem(ctx context.Context, title, body string, private, tweet bool) (*Item, error) {
	// TODO: implement
	return nil, nil
}

// UpdateItem update the item having provided itemID.
// This method requires authentication.
//
// PATCH /api/v2/items/:item_id
// document: http://qiita.com/api/v2/docs#patch-apiv2itemsitem_id
func (c *Client) UpdateItem(ctx context.Context, itemID string, title, body string, private, tweet bool) (*Item, error) {
	// TODO: implement
	return nil, nil
}

// DeleteItem deletes the item.
// This method requires authentication.
//
// DELETE /api/v2/items/:item_id
// document: https://qiita.com/api/v2/docs#delete-apiv2itemsitem_id
func (c *Client) DeleteItem(ctx context.Context, itemID string) error {
	// TODO: implement
	return nil
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
	// TODO: implement
	return false, nil
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
