package service

import "github.com/kolibriee/users-rest-api/internal/entities"

type Authorization interface {
	SignUp(user entities.User) (int, error)
	SignIn(signInUser entities.SignInUserInput) (string, string, error)
	Refresh(refreshToken string) (string, string, error)
}

type Service struct {
	Authorization
}

func NewService() *Service {
	return &Service{
		Authorization: NewAuthorizationService(),
	}
}
