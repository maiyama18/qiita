package main

// Tag represents a tag for a qiita item.
type Tag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}
