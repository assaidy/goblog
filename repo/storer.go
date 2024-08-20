package repo

import "github.com/assaidy/goblog/models"

// Storer defines the interface for user data storage operations.
// Implementations of this interface should provide methods for CRUD operations and checks.
type Storer interface {
	CreateUser(*models.User) (*models.User, error)
	GetUserById(int) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	UpdateUserById(int, *models.UserRegisterOrUpdateRequest) (*models.User, error)
	DeleteUserById(int) error
	GetAllUsers() ([]*models.User, error)
	IsUsernameUsed(string) (bool, error)
	IsEmailUsed(string) (bool, error)
}

