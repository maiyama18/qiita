package qiita

import (
	"context"
)

// Tag represents tag which can be attached to a qiita item.
type Tag struct {
	ID             string `json:"id"`
	IconURL        string `json:"icon_url"`
	ItemsCount     int    `json:"items_count"`
	FollowersCount int    `json:"followers_count"`
}

// GetTag fetches the tag having provided tagID.
//
// GET /api/v2/tags/:tag_id
// document: http://qiita.com/api/v2/docs#get-apiv2tagstag_id
func (c *Client) GetTag(ctx context.Context, tagID string) (*Tag, error) {
	// TODO: implement
	return nil, nil
}

// GetTags fetches all the tags.
//
// GET /api/v2/tags
// document: http://qiita.com/api/v2/docs#get-apiv2tags
func (c *Client) GetTags(ctx context.Context) ([]*Tag, error) {
	// TODO: implement
	return nil, nil
}

// GetTagItems fetches the items which the tag having provided tagID is attached.
//
// GET /api/v2/tags/:tag_id/items
// document: http://qiita.com/api/v2/docs#get-apiv2tagstag_iditems
func (c *Client) GetTagItems(ctx context.Context, tagID string) ([]*Item, error) {
	// TODO: implement
	return nil, nil
}

// IsFollowingTag returns true if the authenticated user is following the tag having provided tagID.
// This method requires authentication.
//
// GET /api/v2/tags/:tag_id/following
// document: http://qiita.com/api/v2/docs#get-apiv2tagstag_idfollowing
func (c *Client) IsFollowingTag(ctx context.Context, tagID string) (bool, error) {
	// TODO: implement
	return false, nil
}

// FollowTag follows the tag having provided tagID.
// This method requires authentication.
//
// PUT /api/v2/tags/:tag_id/following
// document: http://qiita.com/api/v2/docs#put-apiv2tagstag_idfollowing
func (c *Client) FollowTag(ctx context.Context, tagID string) error {
	// TODO: implement
	return nil
}

// UnfollowTag unfollows the tag having provided tagID.
// This method requires authentication.
//
// DELETE /api/v2/tags/:tag_id/following
// document: http://qiita.com/api/v2/docs#delete-apiv2tagstag_idfollowing
func (c *Client) UnfollowTag(ctx context.Context, tagID string) error {
	// TODO: implement
	return nil
}
