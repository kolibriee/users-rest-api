package entities

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type CreateUserInput struct {
	Role     string `json:"role" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	City     string `json:"city" validate:"required"`
}

type SignUpInput struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	City     string `json:"city" validate:"required"`
}

type SignInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateInput struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	City     *string `json:"city"`
	Role     *string `json:"role"`
}

func (input *CreateUserInput) ValidateCreateUserInput() error {
	if err := validate.Struct(input); err != nil {
		return err
	}
	if input.Role != "admin" && input.Role != "user" {
		return errors.New("role must be admin or user")
	}
	return nil
}

func (input *SignUpInput) ValidateSignUpInput() error {
	return validate.Struct(input)
}

func (input *SignInInput) ValidateSignInInput() error {
	return validate.Struct(input)
}

func (u UserUpdateInput) ValidateUserUpdate(role string) error {
	// Проверяем, что хотя бы одно поле для обновления не является nil
	if u.Name == nil && u.Username == nil && u.Password == nil && u.Email == nil && u.City == nil && u.Role == nil {
		return errors.New("update must have at least one of: name, username, password, email, city, or role")
	}

	// Проверяем, что поля не пустые (если они не nil)
	if u.Password != nil {
		if *u.Password == "" {
			return errors.New("password must not be empty")
		}
		if len(*u.Password) < 6 {
			return errors.New("password must be at least 6 characters long")
		}
	}
	if u.Username != nil && *u.Username == "" {
		return errors.New("username must not be empty")
	}
	if u.Email != nil && *u.Email == "" {
		return errors.New("email must not be empty")
	}

	// Только админ может обновлять роль
	if role != "admin" && u.Role != nil {
		return errors.New("only admin can update role")
	}

	return nil
}
