package service

import (
	"context"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/kolibriee/users-rest-api/internal/repository"
	"github.com/kolibriee/users-rest-api/pkg/auth"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) GetAllUsers(ctx context.Context) ([]bunEntities.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UsersService) GetUserByID(ctx context.Context, id int) (*bunEntities.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UsersService) CreateUser(ctx context.Context, user entities.CreateUserInput) (int, error) {
	user.Password = auth.GeneratePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

func (s *UsersService) UpdateUser(ctx context.Context, userID int, user entities.UserUpdateInput) error {
	if user.Password != nil {
		*user.Password = auth.GeneratePasswordHash(*user.Password)
	}
	return s.repo.UpdateUser(ctx, userID, user)
}

func (s *UsersService) DeleteUser(ctx context.Context, userID int) error {
	return s.repo.DeleteUser(ctx, userID)
}
