package models

type Storager interface {
	CreateArticle(*Article) error
	UpdateArticle(int, string, string) error
	GetArticle(int) (*Article, error)
	DeleteArticle(int) error
	GetAllArticles() ([]Article, error)
	DeleteAllArticles() error
}
