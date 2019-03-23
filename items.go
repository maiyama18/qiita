package main

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"
)

// Item represents an article published on qiita
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

	User *User  `json:"user"`
	Tags []*Tag `json:"tags"`
}

// Tag represents a tag for a qiita item
type Tag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

// GetItem fetches an item having provided itemID
//
// GET /api/v2/items/:item_id
// document: https://qiita.com/api/v2/docs#get-apiv2itemsitem_id
func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("items", itemID), nil, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("item with id '%s' not found (status = %d)", itemID, resp.StatusCode)
		default:
			return nil, fmt.Errorf("unknown error (status = %d)", resp.StatusCode)
		}
	}

	var item Item
	if err := c.decodeBody(resp, &item); err != nil {
		return nil, err
	}

	return &item, nil
}
