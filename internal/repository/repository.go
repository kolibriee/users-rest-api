package repository

import (
	"context"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/uptrace/bun"
)

type Authorization interface {
	CreateUser(ctx context.Context, user entities.SignUpInput) (int, error)
	GetUser(ctx context.Context, signinuser entities.SignInInput) (bunEntities.User, error)
	CreateSession(ctx context.Context, session bunEntities.Session) (string, error)
	GetSession(ctx context.Context, refreshToken string) (bunEntities.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
	GetRole(ctx context.Context, userID int) (string, error)
}

type Users interface {
	GetAllUsers(ctx context.Context) ([]bunEntities.User, error)
	GetUserByID(ctx context.Context, id int) (*bunEntities.User, error)
	CreateUser(ctx context.Context, user entities.CreateUserInput) (int, error)
	UpdateUser(ctx context.Context, userID int, user entities.UserUpdateInput) error
	DeleteUser(ctx context.Context, id int) error
}

type Repository struct {
	Authorization
	Users
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Users:         NewUsersRepository(db),
	}
}
