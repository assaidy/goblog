package models

import (
	"strings"
	"time"
)

type Post struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  int       `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostCreateOrUpdateRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId int    `json:"authorId"`
}

// Validate checks if the PostCreateOrUpdateRequest fields are valid.
func (r *PostCreateOrUpdateRequest) Validate() []string {
	var errors []string
	r.Title = strings.TrimSpace(r.Title)
	r.Content = strings.TrimSpace(r.Content)
	if r.Title == "" {
		errors = append(errors, "title is required")
	}
	if r.Content == "" {
		errors = append(errors, "content is required")
	}
	return errors
}
