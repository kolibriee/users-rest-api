package v1

import (
	"net/http"

	"github.com/kolibriee/users-rest-api/internal/entities"
	"github.com/labstack/echo/v4"
)

// SignUp godoc
//
//	@Summary	Register a new user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		entities.SignUpInput	true	"SignUp input"
//	@Success	200		{object}	map[string]interface{}
//	@Failure	400		{object}	string	"Invalid input body"
//	@Failure	500		{object}	string	"Internal server error"
//	@Router		/auth/sign-up [post]
func (h *Handler) SignUp(c echo.Context) error {
	var input entities.SignUpInput
	if err := c.Bind(&input); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid input body") // Обработка ошибки привязки
	}

	if err := input.ValidateSignUpInput(); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error()) // Проверка на валидность
	}
	id, err := h.services.Authorization.SignUp(c.Request().Context(), input) // Регистрация пользователя
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// SignIn godoc
//
//	@Summary	Login user and get tokens
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		entities.SignInInput	true	"SignIn input"
//	@Success	200		{object}	map[string]interface{}
//	@Failure	400		{object}	string	"Invalid input body"
//	@Failure	500		{object}	string	"Internal server error"
//	@Router		/auth/sign-in [post]
func (h *Handler) SignIn(c echo.Context) error {
	var input entities.SignInInput
	if err := c.Bind(&input); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid input body") // Обработка ошибки привязки
	}

	if err := input.ValidateSignInInput(); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error()) // Проверка на валидность
	}

	accessToken, refreshToken, err := h.services.Authorization.SignIn(c.Request().Context(), input) // Авторизация пользователя
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error()) // Обработка ошибки
	}

	// Установка куки с refreshToken
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}

// Refresh godoc
//
//	@Summary	Refresh access token
//	@Tags		auth
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Failure	401	{object}	string	"No refresh token provided"
//	@Failure	500	{object}	string	"Internal server error"
//	@Router		/auth/refresh [get]
func (h *Handler) Refresh(c echo.Context) error {
	refreshTokenCookie, err := c.Cookie("refreshToken") // Получение refreshToken из куки
	if err != nil {
		return newErrorResponse(c, http.StatusUnauthorized, "no refresh token provided") // Проверка наличия куки
	}

	accessToken, newRefreshToken, err := h.services.Authorization.Refresh(c.Request().Context(), refreshTokenCookie.Value) // Обновление токенов
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error()) // Обработка ошибки
	}

	// Установка нового refreshToken
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}
