package postgres_repo

// TODO: when removing users, use their ID's for new users.

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/assaidy/goblog/models"
	"github.com/assaidy/goblog/utils"
	_ "github.com/lib/pq"
)

// PostgresRepo implements the Storer interface for PostgreSQL.
type PostgresRepo struct {
	DB *sql.DB
}

// NewPostgresRepo initializes a new PostgresRepo with the given database connection string.
func NewPostgresRepo(dbConn string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepo{DB: db}, nil
}

// CreateUser inserts a new user into the database and returns the created user.
func (pg *PostgresRepo) CreateUser(user *models.User) (*models.User, error) {
	query := `
    INSERT INTO users (full_name, username, email, password, bio, joined_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id;`

	err := pg.DB.QueryRow(query, user.FullName, user.Username, user.Email, user.Password, user.Bio, user.JoinedAt).Scan(&user.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserById retrieves a user by their ID.
func (pg *PostgresRepo) GetUserById(id int) (*models.User, error) {
	query := `
    SELECT id, full_name, username, email, password, bio, joined_at
    FROM users
    WHERE id = $1`

	user := &models.User{}
	err := pg.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.JoinedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NotFound(fmt.Errorf("no user with id %d", id))
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func (pg *PostgresRepo) GetUserByUsername(username string) (*models.User, error) {
	query := `
    SELECT id, full_name, username, email, password, bio, joined_at
    FROM users
    WHERE username = $1`

	user := &models.User{}
	err := pg.DB.QueryRow(query, username).Scan(
		&user.Id,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.JoinedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NotFound(fmt.Errorf("no user with username %s", username))
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database.
func (pg *PostgresRepo) GetAllUsers() ([]*models.User, error) {
	query := `
    SELECT id, full_name, username, email, bio, joined_at
    FROM users`

	rows, err := pg.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Username,
			&user.Email,
			&user.Bio,
			&user.JoinedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserById updates an existing user identified by ID with new information.
func (pg *PostgresRepo) UpdateUserById(id int, updateReq *models.UserRegisterOrUpdateRequest) (*models.User, error) {
	query := `
    UPDATE users SET
        full_name = $1,
        username = $2,
        email = $3,
        password = $4,
        bio = $5
    WHERE id = $6
    RETURNING joined_at`

	user := &models.User{
		Id:       id,
		FullName: updateReq.FullName,
		Username: updateReq.Username,
		Email:    updateReq.Email,
		Password: updateReq.Password,
		Bio:      updateReq.Bio,
	}

	err := pg.DB.QueryRow(query,
		updateReq.FullName,
		updateReq.Username,
		updateReq.Email,
		updateReq.Password,
		updateReq.Bio,
		id,
	).Scan(&user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFound(fmt.Errorf("no user with id %d", id))
		}
		return nil, err
	}

	return user, nil
}

// DeleteUserById removes a user identified by ID from the database.
func (pg *PostgresRepo) DeleteUserById(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := pg.DB.Exec(query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFound(fmt.Errorf("no user with id %d", id))
	}

	return nil
}

// IsUsernameUsed checks if the provided username is already in use.
func (pg *PostgresRepo) IsUsernameUsed(username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = $1 LIMIT 1`

	var exists int
	err := pg.DB.QueryRow(query, username).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exists == 1, nil
}

// IsEmailUsed checks if the provided email is already in use.
func (pg *PostgresRepo) IsEmailUsed(email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`

	var exists int
	err := pg.DB.QueryRow(query, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exists == 1, nil
}

func (pg *PostgresRepo) CreatePost(post *models.Post) (*models.Post, error) {
	query := `
    INSERT INTO posts (title, content, author_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id;`

	err := pg.DB.QueryRow(query, post.Title, post.Content, post.AuthorId, post.CreatedAt, post.UpdatedAt).Scan(&post.Id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresRepo) GetPostById(id int) (*models.Post, error) {
	query := `
    SELECT id, title, content, author_id, created_at, updated_at 
    FROM posts
    WHERE id = $1`

	post := &models.Post{}
	err := pg.DB.QueryRow(query, id).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.AuthorId,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NotFound(fmt.Errorf("no post with id %d", id))
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresRepo) UpdatePostById(id int, postReq *models.PostCreateOrUpdateRequest) (*models.Post, error) {
	query := `
    UPDATE posts SET
        title = $1,
        content = $2,
        updated_at = $3
    WHERE id = $4
    RETURNING author_id, created_at`

	post := &models.Post{
		Id:        id,
		Title:     postReq.Title,
		Content:   postReq.Content,
		UpdatedAt: time.Now().UTC(),
	}

	err := pg.DB.QueryRow(query, post.Title, post.Content, post.UpdatedAt, post.Id).Scan(
		&post.AuthorId,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFound(fmt.Errorf("no post with id %d", id))
		}
		return nil, err
	}

	return post, nil
}

func (pg *PostgresRepo) DeletePostById(id, authorId int) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := pg.DB.Exec(query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFound(fmt.Errorf("no post with id %d found", id))
	}

	return nil
}

func (pg *PostgresRepo) GetAllPosts() ([]*models.Post, error) {
	query := `
    SELECT id, title, content, author_id, created_at, updated_at 
    FROM posts`

	rows, err := pg.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (pg *PostgresRepo) GetAllPostsByAuthor(authorId int) ([]*models.Post, error) {
	query := `
    SELECT id, title, content, author_id, created_at, updated_at 
    FROM posts
    WHERE author_id = $1`

	rows, err := pg.DB.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
