package models

import "time"

type Article struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishDate time.Time `json:"publishDate"`
}

type ArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
