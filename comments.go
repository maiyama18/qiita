package qiita

import "context"

// Comment represents a comment on qiita item.
type Comment struct{}

// GetComment fetches the comment having provided commentID.
//
// GET /api/v2/comments/:comment_id
// document: http://qiita.com/api/v2/docs#get-apiv2commentscomment_id
func (c *Client) GetComment(ctx context.Context, commentID string) (*Comment, error) {
	// TODO: implement
	return nil, nil
}

// UpdateComment update the comment having provided commentID.
// This method requires authentication.
//
// PATCH /api/v2/comments/:comment_id
// document: https://qiita.com/api/v2/docs#patch-apiv2commentscomment_id
func (c *Client) UpdateComment(ctx context.Context, commentID string, body string) error {
	// TODO: implement
	return nil
}

// DeleteComment delete the comment having provided commentID.
// This method requires authentication.
//
// DELETE /api/v2/comments/:comment_id
// document: http://qiita.com/api/v2/docs#delete-apiv2commentscomment_id
func (c *Client) DeleteComment(ctx context.Context, commentID string) error {
	// TODO: implement
	return nil
}

// ThankComment post thank on the comment having provided commentID.
// This method requires authentication.
//
// PUT /api/v2/comments/:comment_id/thank
// document: http://qiita.com/api/v2/docs#put-apiv2commentscomment_idthank
func (c *Client) ThankComment(ctx context.Context, commentID string) error {
	// TODO: implement
	return nil
}

// UnthankComment delete thank on the comment having provided commentID.
// This method requires authentication.
//
// DELETE /api/v2/comments/:comment_id/thank
// document: http://qiita.com/api/v2/docs#delete-apiv2commentscomment_idthank
func (c *Client) UnthankComment(ctx context.Context, commentID string) error {
	// TODO: implement
	return nil
}
