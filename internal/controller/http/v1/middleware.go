package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kolibriee/users-rest-api/pkg/auth"
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	adminCtx            = "adminId"
	roleCtx             = "role"
)

// UserIdentity middleware для проверки идентификации пользователя
func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeader) // Получение заголовка авторизации
		if authHeader == "" {
			return newErrorResponse(c, http.StatusUnauthorized, "empty auth header") // Проверка наличия заголовка
		}
		headerParts := strings.Split(authHeader, " ") // Разделение заголовка на части
		if len(headerParts) != 2 {
			return newErrorResponse(c, http.StatusUnauthorized, "invalid auth header") // Проверка формата заголовка
		}
		userId, role, err := auth.ParseToken(headerParts[1]) // Парсинг токена
		if err != nil {
			return newErrorResponse(c, http.StatusUnauthorized, err.Error()) // Обработка ошибки парсинга
		}
		c.Set(userCtx, userId) // Установка идентификатора пользователя в контекст
		c.Set(roleCtx, role)   // Установка роли в контекст
		return next(c)         // Передача управления следующему обработчику
	}
}

// AdminIdentity middleware для проверки идентификации администратора
func (h *Handler) adminIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeader) // Получение заголовка авторизации
		if authHeader == "" {
			return newErrorResponse(c, http.StatusUnauthorized, "empty auth header") // Проверка наличия заголовка
		}
		headerParts := strings.Split(authHeader, " ") // Разделение заголовка на части
		if len(headerParts) != 2 {
			return newErrorResponse(c, http.StatusUnauthorized, "invalid auth header") // Проверка формата заголовка
		}
		userId, role, err := auth.ParseToken(headerParts[1]) // Парсинг токена
		if err != nil {
			return newErrorResponse(c, http.StatusUnauthorized, err.Error()) // Обработка ошибки парсинга
		}
		if role != "admin" { // Проверка роли
			return newErrorResponse(c, http.StatusForbidden, "access denied") // Отказ в доступе
		}
		c.Set(userCtx, userId) // Установка идентификатора пользователя в контекст
		c.Set(roleCtx, role)   // Установка роли в контекст
		return next(c)         // Передача управления следующему обработчику
	}
}

// getUserId извлекает userId из контекста
func getUserId(c echo.Context) (int, error) {
	id, ok := c.Get(userCtx).(int) // Извлечение идентификатора пользователя
	if !ok {
		return 0, errors.New("user id not found")
	}
	return id, nil
}

// getAdminId извлекает adminId из контекста
func getAdminId(c echo.Context) (int, error) {
	id, ok := c.Get(adminCtx).(int) // Извлечение идентификатора администратора
	if !ok {
		return 0, errors.New("admin id not found")
	}
	return id, nil
}

// getRole извлекает role из контекста
func getRole(c echo.Context) (string, error) {
	role, ok := c.Get(roleCtx).(string) // Извлечение роли
	if !ok {
		return "", errors.New("role not found")
	}
	return role, nil
}
