package service

import (
	"context"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/kolibriee/users-rest-api/internal/repository"
)

type Authorization interface {
	SignUp(ctx context.Context, user entities.SignUpInput) (int, error)
	SignIn(ctx context.Context, ignInUser entities.SignInInput) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type Users interface {
	GetAllUsers(ctx context.Context) ([]bunEntities.User, error)
	GetUserByID(ctx context.Context, id int) (*bunEntities.User, error)
	CreateUser(ctx context.Context, user entities.CreateUserInput) (int, error)
	UpdateUser(ctx context.Context, userID int, user entities.UserUpdateInput) error
	DeleteUser(ctx context.Context, id int) error
}

type Service struct {
	Authorization
	Users
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repo.Authorization),
		Users:         NewUsersService(repo.Users),
	}
}
