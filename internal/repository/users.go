package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kolibriee/users-rest-api/internal/entities"
	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/uptrace/bun"
)

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetAllUsers(ctx context.Context) ([]bunEntities.User, error) {
	var users []bunEntities.User
	if err := r.db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UsersRepository) GetUserByID(ctx context.Context, id int) (*bunEntities.User, error) {
	var user bunEntities.User
	if err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepository) CreateUser(ctx context.Context, user entities.CreateUserInput) (int, error) {
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
		Role:         user.Role,
		Name:         user.Name,
		Username:     user.Username,
		PasswordHash: user.Password,
		City:         user.City,
	}

	// Сохраняем пользователя в БД
	_, err = r.db.NewInsert().Model(newUser).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return newUser.ID, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, userID int, user entities.UserUpdateInput) error {
	// Если обновляется имя пользователя, проверяем, существует ли пользователь с таким же именем
	if user.Username != nil {
		existingUser := &bunEntities.User{}
		err := r.db.NewSelect().Model(existingUser).
			Where("username = ?", *user.Username).
			Where("id != ?", userID). // Исключаем текущего пользователя из поиска
			Limit(1).
			Scan(ctx)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if existingUser.ID > 0 {
			return errors.New("user with this username already exists")
		}
	}

	// Подготовка обновляемых полей
	updatedUser := &bunEntities.User{}
	columnsToUpdate := []string{}

	if user.Role != nil {
		updatedUser.Role = *user.Role
		columnsToUpdate = append(columnsToUpdate, "role")
	}
	if user.Name != nil {
		updatedUser.Name = *user.Name
		columnsToUpdate = append(columnsToUpdate, "name")
	}
	if user.Username != nil {
		updatedUser.Username = *user.Username
		columnsToUpdate = append(columnsToUpdate, "username")
	}
	if user.Password != nil {
		updatedUser.PasswordHash = *user.Password
		columnsToUpdate = append(columnsToUpdate, "password_hash")
	}
	if user.City != nil {
		updatedUser.City = *user.City
		columnsToUpdate = append(columnsToUpdate, "city")
	}

	// Выполняем обновление только тех полей, которые нужно изменить
	_, err := r.db.NewUpdate().
		Model(updatedUser).
		Column(columnsToUpdate...).
		Where("id = ?", userID).
		Exec(ctx)

	return err
}

func (r *UsersRepository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.db.NewDelete().Model((*bunEntities.User)(nil)).Where("id = ?", userID).Exec(ctx)
	return err
}
