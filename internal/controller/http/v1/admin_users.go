package v1

import (
	"errors"
	"net/http"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bun_entities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/labstack/echo/v4"
)

var _ = bun_entities.User{}

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Get list of all users (admin only)
//	@Tags			admin
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{array}		bun_entities.User
//	@Failure		403	{object}	statusResponse	"access denied"
//	@Failure		500	{object}	statusResponse	"internal server error"
//	@Router			/admin/users [get]
func (h *Handler) GetAllUsers(c echo.Context) error {
	role, err := getRole(c) // Получение роли пользователя
	if err != nil || role != "admin" {
		return newErrorResponse(c, http.StatusForbidden, errors.New("access denied").Error()) // Проверка на доступ
	}

	users, err := h.services.GetAllUsers(c.Request().Context()) // Получение всех пользователей
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, errors.New("can't get users; ").Error()+err.Error()) // Обработка ошибки
	}

	return c.JSON(http.StatusOK, users) // Возвращение списка пользователей
}

// GetUserByID godoc
//
//	@Summary		Get user by ID
//	@Description	Get a user by their ID (admin or user themselves)
//	@Tags			admin
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	bun_entities.User
//	@Failure		400	{object}	statusResponse	"invalid user id"
//	@Failure		403	{object}	statusResponse	"access denied"
//	@Failure		500	{object}	statusResponse	"internal server error"
//	@Router			/admin/users/{id} [get]
func (h *Handler) AdminGetUserByID(c echo.Context) error {
	return h.GetUserByID(c)
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user (admin only)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			user	body		entities.CreateUserInput	true	"User data"
//	@Success		201		{object}	map[string]interface{}		"user created"
//	@Failure		400		{object}	statusResponse				"invalid request"
//	@Failure		500		{object}	statusResponse				"internal server error"
//	@Router			/admin/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var user entities.CreateUserInput
	if err := c.Bind(&user); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, errors.New("invalid request").Error())
	}
	if err := user.ValidateCreateUserInput(); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error()) // Проверка на валидность данных
	}
	id, err := h.services.CreateUser(c.Request().Context(), user) // Создание нового пользователя
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, errors.New("can't create user; ").Error()+err.Error()) // Обработка ошибки
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update user information (admin or user themselves)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int							true	"User ID"
//	@Param			user	body		entities.UserUpdateInput	true	"Updated user data"
//	@Success		200		{object}	statusResponse				"ok"
//	@Failure		400		{object}	statusResponse				"invalid request"
//	@Failure		403		{object}	statusResponse				"access denied"
//	@Failure		500		{object}	statusResponse				"internal server error"
//	@Router			/admin/users/{id} [put]
func (h *Handler) AdminUpdateUser(c echo.Context) error {
	return h.UpdateUser(c)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by their ID (admin or user themselves)
//	@Tags			admin
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int				true	"User ID"
//	@Success		200	{object}	statusResponse	"ok"
//	@Failure		400	{object}	statusResponse	"invalid user id"
//	@Failure		403	{object}	statusResponse	"access denied"
//	@Failure		500	{object}	statusResponse	"internal server error"
//	@Router			/admin/users/{id} [delete]
func (h *Handler) AdminDeleteUser(c echo.Context) error {
	return h.DeleteUser(c)
}
