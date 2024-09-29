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
)

// UserIdentity middleware для проверки идентификации пользователя
func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeader)
		if authHeader == "" {
			return newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		}
		userId, _, err := auth.ParseToken(headerParts[1])
		if err != nil {
			return newErrorResponse(c, http.StatusUnauthorized, err.Error())
		}
		c.Set(userCtx, userId)
		return next(c) // передаем управление следующему обработчику
	}
}

// AdminIdentity middleware для проверки идентификации администратора
func (h *Handler) adminIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeader)
		if authHeader == "" {
			return newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		}
		userId, role, err := auth.ParseToken(headerParts[1])
		if err != nil {
			return newErrorResponse(c, http.StatusUnauthorized, err.Error())
		}
		if role != "admin" { // проверяем роль
			return newErrorResponse(c, http.StatusForbidden, "access denied")
		}
		c.Set(userCtx, userId)
		return next(c)
	}
}

// getUserId извлекает userId из контекста
func getUserId(c echo.Context) (int, error) {
	id, ok := c.Get(userCtx).(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return id, nil
}

// getAdminId извлекает adminId из контекста
func getAdminId(c echo.Context) (int, error) {
	id, ok := c.Get(adminCtx).(int)
	if !ok {
		return 0, errors.New("admin id not found")
	}
	return id, nil
}
