package main

import "time"

// TODO: add tags, user, group
type Item struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Body         string `json:"body"`
	RenderedBody string `json:"rendered_body"`
	Private      bool   `json:"private"`
	Coediting    bool   `json:"coediting"`

	CreatedAt time.Time `json:"created_at""`
	UpdatedAt time.Time `json:"updated_at""`

	PageViewsCount int `json:"page_views_count"`
	CommentsCount  int `json:"comments_count"`
	LikesCount     int `json:"likes_count"`
	ReactionsCount int `json:"reactions_count"`
}
