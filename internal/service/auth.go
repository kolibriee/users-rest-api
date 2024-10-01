package service

import (
	"context"
	"errors"
	"time"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/kolibriee/users-rest-api/internal/repository"
	"github.com/kolibriee/users-rest-api/pkg/auth"
)

const (
	accessTokenTTL  = 30 * time.Minute    // Время жизни access token
	refreshTokenTTL = 30 * 24 * time.Hour // Время жизни refresh token
)

type AuthorizationService struct {
	repo repository.Authorization
}

// NewAuthorizationService создает новый экземпляр AuthorizationService с переданным репозиторием.
func NewAuthorizationService(repo repository.Authorization) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

// SignUp регистрирует нового пользователя и возвращает его ID.
func (s *AuthorizationService) SignUp(ctx context.Context, user entities.SignUpInput) (int, error) {
	user.Password = auth.GeneratePasswordHash(user.Password) // Хешируем пароль
	return s.repo.CreateUser(ctx, user)
}

// SignIn выполняет вход пользователя, возвращая access token и refresh token.
func (s *AuthorizationService) SignIn(ctx context.Context, signInUser entities.SignInInput) (string, string, error) {
	signInUser.Password = auth.GeneratePasswordHash(signInUser.Password) // Хешируем пароль
	user, err := s.repo.GetUser(ctx, signInUser)                         // Получаем пользователя из репозитория
	if err != nil {
		return "", "", err // Возвращаем ошибку, если пользователь не найден
	}

	// Генерируем access token
	accessToken, err := auth.GenerateAccessToken(accessTokenTTL, user.ID, user.Role)
	if err != nil {
		return "", "", err // Возвращаем ошибку, если не удалось сгенерировать token
	}

	// Создаем новую сессию для пользователя
	refreshToken, err := s.repo.CreateSession(ctx, bunEntities.Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenTTL), // Устанавливаем время истечения сессии
	})
	if err != nil {
		return "", "", errors.New("can't create refresh token" + err.Error())
	}
	return accessToken, refreshToken, nil
}

// Refresh обновляет access token и refresh token, используя refresh token.
func (s *AuthorizationService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.repo.GetSession(ctx, refreshToken) // Получаем сессию по refresh token
	if err != nil {
		return "", "", errors.New("invalid refresh token" + err.Error())
	}

	if session.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("session expired") // Проверяем, истекла ли сессия
	}

	role, err := s.repo.GetRole(ctx, session.UserID) // Получаем роль пользователя
	accessToken, err := auth.GenerateAccessToken(accessTokenTTL, session.UserID, role)
	if err != nil {
		return "", "", errors.New("can't generate access token" + err.Error())
	}

	// Создаем новую сессию с новым refresh token
	newRefreshToken, err := s.repo.CreateSession(ctx, bunEntities.Session{
		UserID:    session.UserID,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	})
	if err != nil {
		return "", "", errors.New("can't create refresh token" + err.Error())
	}

	// Удаляем старую сессию
	if err := s.repo.DeleteSession(ctx, refreshToken); err != nil {
		return "", "", errors.New("can't delete old refresh token" + err.Error())
	}
	return accessToken, newRefreshToken, nil // Возвращаем новые токены
}
