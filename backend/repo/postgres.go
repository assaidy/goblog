package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/assaidy/personal-blog/backend/models"
	_ "github.com/lib/pq"
)

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (models.Storager, error) {
	connStr := "user=postgres dbname=goblog password=goblog sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, models.NewDBError(err)
	}

	if err := db.Ping(); err != nil {
		return nil, models.NewDBError(err)
	}

	pg := &postgresStore{db: db}

	if err := pg.Init(); err != nil {
		log.Fatal(err)
	}

	return pg, nil
}

func (pg *postgresStore) Init() error {
	return pg.createArticleTable()
}

func (pg *postgresStore) createArticleTable() error {
	query := `CREATE TABLE IF NOT EXISTS article (
        id SERIAL PRIMARY KEY,
        title VARCHAR(150),
        content TEXT,
        publish_date TIMESTAMP
    )`

	_, err := pg.db.Exec(query)
	return err
}

func (pg *postgresStore) CreateArticle(a *models.Article) error {
	query := `INSERT INTO article (title, content, publish_date)
              VALUES ($1, $2, $3)`

	_, err := pg.db.Exec(query, a.Title, a.Content, a.PublishDate)
	if err != nil {
		return models.NewDBError(err)
	}

	return nil
}

func (pg *postgresStore) UpdateArticle(id int, newTitle, newContent string) error {
	query := `UPDATE article SET title = $1, content = $2 WHERE id = $3`

	result, err := pg.db.Exec(query, newTitle, newContent, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewDBError(fmt.Errorf("error checking rows affected: %w", err))
	}

	if rowsAffected == 0 {
		return models.NewNotFoundError("article", id)
	}

	return nil
}

func (pg *postgresStore) GetArticle(id int) (*models.Article, error) {
	query := `SELECT id, title, content, publish_date FROM article WHERE id = $1`

	var article models.Article

	err := pg.db.QueryRow(query, id).Scan(&article.Id, &article.Title, &article.Content, &article.PublishDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewNotFoundError("article", id)
		}
		return nil, models.NewDBError(err)
	}

	return &article, nil
}

func (pg *postgresStore) DeleteArticle(id int) error {
	query := `DELETE FROM article WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return models.NewDBError(fmt.Errorf("error deleting article with id %d: %w", id, err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.NewDBError(fmt.Errorf("error checking rows affected: %w", err))
	}

	if rowsAffected == 0 {
		return models.NewNotFoundError("article", id)
	}

	return nil
}

func (pg *postgresStore) GetAllArticles() ([]models.Article, error) {
	query := `SELECT id, title, content, publish_date FROM article`

	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, models.NewDBError(err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0)

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.Title, &article.Content, &article.PublishDate)
		if err != nil {
			return nil, models.NewDBError(fmt.Errorf("error scanning row: %w", err))
		}

		articles = append(articles, article)
	}

	// Check for errors after the loop has completed.
	if err = rows.Err(); err != nil {
		return nil, models.NewDBError(fmt.Errorf("error iterating over rows: %w", err))
	}

	return articles, nil
}

func (pg *postgresStore) DeleteAllArticles() error {
	query := `DELETE FROM article`

	_, err := pg.db.Exec(query)
	if err != nil {
		return models.NewDBError(fmt.Errorf("error deleting article: %w", err))
	}

	return nil
}
