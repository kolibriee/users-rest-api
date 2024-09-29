package service

import "github.com/kolibriee/users-rest-api/internal/entities"

type AuthorizationService struct {
}

func NewAuthorizationService() *AuthorizationService {
	return &AuthorizationService{}
}

func (s *AuthorizationService) SignUp(user entities.User) (int, error) {

	return 0, nil
}

func (s *AuthorizationService) SignIn(signInUser entities.SignInUserInput) (string, string, error) {

	return "", "", nil
}

func (s *AuthorizationService) Refresh(refreshToken string) (string, string, error) {

	return "", "", nil
}
