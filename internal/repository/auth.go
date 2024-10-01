package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
)

type AuthRepository struct {
	db *bun.DB // База данных для выполнения операций
}

func NewAuthRepository(db *bun.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных.
func (r *AuthRepository) CreateUser(ctx context.Context, user entities.SignUpInput) (int, error) {
	// Проверяем, существует ли пользователь с таким же именем пользователя
	existingUser := &bunEntities.User{}
	err := r.db.NewSelect().Model(existingUser).
		Where("username = ?", user.Username).
		Limit(1).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	// Если пользователь найден (ID больше 0), возвращаем ошибку
	if existingUser.ID > 0 {
		return 0, errors.New("user with this username already exists")
	}

	// Создаем нового пользователя
	newUser := &bunEntities.User{
		Role:         "user", // Роль по умолчанию
		Name:         user.Name,
		Username:     user.Username,
		PasswordHash: user.Password,
		City:         user.City,
	}

	// Сохраняем пользователя в БД
	_, err = r.db.NewInsert().Model(newUser).Exec(ctx)
	if err != nil {
		return 0, err // Возвращаем ошибку, если не удалось сохранить пользователя
	}

	return newUser.ID, nil
}

// GetUser возвращает пользователя по его имени пользователя и паролю.
func (r *AuthRepository) GetUser(ctx context.Context, signinuser entities.SignInInput) (bunEntities.User, error) {
	var user bunEntities.User

	// Выполняем выборку пользователя по username и password_hash
	err := r.db.NewSelect().
		Model(&user).
		Column("id", "role", "name", "username").
		Where("username = ?", signinuser.Username).
		Where("password_hash = ?", signinuser.Password).
		Scan(ctx)

	if err != nil {
		return bunEntities.User{}, err // Возвращаем ошибку, если не удалось найти пользователя
	}

	return user, nil
}

// CreateSession создает новую сессию и возвращает сгенерированный refresh_token.
func (r *AuthRepository) CreateSession(ctx context.Context, session bunEntities.Session) (string, error) {
	// Вставляем новую сессию в базу данных, исключая refresh_token (он будет сгенерирован базой)
	_, err := r.db.NewInsert().
		Model(&session).
		ExcludeColumn("refresh_token"). // Исключаем refresh_token из запроса
		Returning("refresh_token").     // Возвращаем сгенерированный refresh_token
		Exec(ctx)

	if err != nil {
		return "", err // Возвращаем ошибку, если не удалось создать сессию
	}

	return session.RefreshToken, nil
}

// GetSession возвращает сессию по refresh_token.
func (r *AuthRepository) GetSession(ctx context.Context, refreshToken string) (bunEntities.Session, error) {
	var session bunEntities.Session

	// Выполняем выборку сессии по refresh_token
	err := r.db.NewSelect().
		Model(&session).
		Where("refresh_token = ?", refreshToken).
		Scan(ctx)

	if err != nil {
		return bunEntities.Session{}, err // Возвращаем ошибку, если не удалось найти сессию
	}

	return session, nil
}

// DeleteSession удаляет сессию по refresh_token.
func (r *AuthRepository) DeleteSession(ctx context.Context, refreshToken string) error {
	// Удаляем сессию по refresh_token
	_, err := r.db.NewDelete().
		Model((*bunEntities.Session)(nil)).
		Where("refresh_token = ?", refreshToken).
		Exec(ctx)

	return err
}

// GetRole возвращает роль пользователя по его ID.
func (r *AuthRepository) GetRole(ctx context.Context, userID int) (string, error) {
	var user bunEntities.User
	// Выполняем выборку роли пользователя по его ID
	err := r.db.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Column("role").
		Scan(ctx)

	if err != nil {
		return "", err
	}

	return user.Role, nil
}
