package v1

import (
	"net/http"

	"github.com/kolibriee/users-rest-api/internal/entities"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SignUp(c echo.Context) error {
	var input entities.User
	if err := c.Bind(&input); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}
	id, err := h.services.Authorization.SignUp(input)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) SignIn(c echo.Context) error {
	var input entities.SignInUserInput
	if err := c.Bind(&input); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}

	accessToken, refreshToken, err := h.services.Authorization.SignIn(input)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Задайте true для HTTPS
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(c, http.StatusUnauthorized, "no refresh token provided")
	}

	accessToken, newRefreshToken, err := h.services.Authorization.Refresh(refreshTokenCookie.Value)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Задайте true для HTTPS
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}
