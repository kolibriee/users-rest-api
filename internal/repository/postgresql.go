package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/kolibriee/users-rest-api/internal/config"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	todoItemsTable  = "todo_items"
	usersListsTable = "users_lists"
	listsItemsTable = "lists_items"
	sessionsTable   = "sessions"
)

func NewPostgresDB(cfg *config.Postgres) (*bun.DB, error) {
	// Создаем подключение через pgdriver
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to ping db: " + err.Error())
	}

	return db, nil
}
