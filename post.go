package main

type Post struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	LikesCount int    `json:"likes_count"`
}
