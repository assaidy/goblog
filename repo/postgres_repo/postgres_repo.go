package postgres_repo

import (
	"database/sql"
	"errors"
	"fmt"
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

	var userId int
	err := pg.DB.QueryRow(query, user.FullName, user.Username, user.Email, user.Password, user.Bio, user.JoinedAt).Scan(&userId)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:       userId,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Bio:      user.Bio,
		JoinedAt: user.JoinedAt,
	}, nil
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
