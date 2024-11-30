package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bun_entities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/labstack/echo/v4"
)

var _ = bun_entities.User{}

// GetUserByID godoc
//
//	@Summary		Get user by ID
//	@Description	Get a user by their ID (admin or user themselves)
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	bun_entities.User
//	@Failure		400	{object}	statusResponse	"invalid user id"
//	@Failure		403	{object}	statusResponse	"access denied"
//	@Failure		500	{object}	statusResponse	"internal server error"
//	@Router			/api/users/{id} [get]
func (h *Handler) GetUserByID(c echo.Context) error {
	// Извлечение ID пользователя из параметров запроса
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(c, http.StatusBadRequest, errors.New("invalid user id").Error())
	}

	// Извлечение текущего ID пользователя из контекста
	currentUserId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Извлечение роли из контекста
	role, err := getRole(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Проверка прав доступа
	if role != "admin" && currentUserId != userId {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Получение пользователя по ID
	user, err := h.services.GetUserByID(c.Request().Context(), userId)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, errors.New("can't get user; ").Error()+err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update user information (admin or user themselves)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int							true	"User ID"
//	@Param			user	body		entities.UserUpdateInput	true	"Updated user data"
//	@Success		200		{object}	statusResponse				"ok"
//	@Failure		400		{object}	statusResponse				"invalid request"
//	@Failure		403		{object}	statusResponse				"access denied"
//	@Failure		500		{object}	statusResponse				"internal server error"
//	@Router			/api/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	// Извлечение ID пользователя из параметров запроса
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(c, http.StatusBadRequest, errors.New("invalid user id").Error())
	}

	// Извлечение текущего ID пользователя из контекста
	currentUserId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Извлечение роли из контекста
	role, err := getRole(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Проверка прав доступа
	if role != "admin" && currentUserId != userId {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	var user entities.UserUpdateInput
	if err := c.Bind(&user); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, errors.New("invalid request").Error())
	}

	// Валидация данных обновления пользователя
	if err := user.ValidateUserUpdate(role); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.UpdateUser(c.Request().Context(), userId, user); err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, errors.New("can't update user; ").Error()+err.Error())
	}

	return c.JSON(http.StatusOK, statusResponse{
		Status: "ok", // Возврат статуса обновления
	})
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by their ID (admin or user themselves)
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int				true	"User ID"
//	@Success		200	{object}	statusResponse	"ok"
//	@Failure		400	{object}	statusResponse	"invalid user id"
//	@Failure		403	{object}	statusResponse	"access denied"
//	@Failure		500	{object}	statusResponse	"internal server error"
//	@Router			/api/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	// Извлечение ID пользователя из параметров запроса
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(c, http.StatusBadRequest, errors.New("invalid user id").Error())
	}

	// Извлечение текущего ID пользователя из контекста
	currentUserId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Извлечение роли из контекста
	role, err := getRole(c)
	if err != nil {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Проверка прав доступа
	if role != "admin" && currentUserId != userId {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error())
	}

	// Удаление пользователя из сервиса
	if err := h.services.DeleteUser(c.Request().Context(), userId); err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, errors.New("can't delete user; ").Error()+err.Error())
	}

	return c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}