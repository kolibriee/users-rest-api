package app

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kolibriee/users-rest-api/internal/config"
)

func Migrate(dbConfig *config.Postgres) error {
	m, err := migrate.New(
		"file://migrations/postgresql",
		"postgres://"+dbConfig.Username+":"+dbConfig.Password+"@"+dbConfig.Host+":"+dbConfig.Port+"/"+dbConfig.DBName+"?sslmode="+dbConfig.SSLMode)
	if err != nil {
		return err
	}
	defer m.Close()

	// Выполните миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
