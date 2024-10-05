package app

import (
	"context"
	"errors"
	"os"

	"github.com/kolibriee/users-rest-api/docs"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"

	bunEntities "github.com/kolibriee/users-rest-api/internal/entities/bun"
	"github.com/kolibriee/users-rest-api/pkg/auth"
)

func initAdmin(db *bun.DB) error {
	ctx := context.Background()

	// Проверяем, существует ли уже администратор
	exists, err := db.NewSelect().
		Model((*bunEntities.User)(nil)).
		Where("role = ?", "admin").
		Exists(ctx)
	if err != nil {
		// Возвращаем ошибку, если не удалось проверить существование администратора
		return errors.New("failed to check if admin exists: " + err.Error())
	}

	if exists {
		// Логируем информацию, если администратор уже существует
		logrus.Info("Admin already exists")
		return nil
	}

	// Создаем нового администратора
	admin := &bunEntities.User{
		Role:         "admin",
		Name:         "admin",
		Username:     "admin",
		PasswordHash: auth.GeneratePasswordHash("admin"), // Генерируем хэш пароля
		City:         "city",
	}

	// Вставляем нового администратора в базу данных
	_, err = db.NewInsert().Model(admin).Exec(ctx)
	if err != nil {
		return errors.New("failed to create admin: " + err.Error())
	}

	logrus.Info("Admin created")
	return nil
}

func initSwagger() {
	host := os.Getenv("SWAGGER_HOST")

	if host == "" {
		host = "localhost:8080"
	}
	docs.SwaggerInfo.Host = host
}
