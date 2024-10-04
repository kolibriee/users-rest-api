package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kolibriee/users-rest-api/internal/config"
	ctrl "github.com/kolibriee/users-rest-api/internal/controller"
	"github.com/kolibriee/users-rest-api/internal/repository"
	"github.com/kolibriee/users-rest-api/internal/server"
	"github.com/kolibriee/users-rest-api/internal/service"
	"github.com/sirupsen/logrus"
)

func Run(configPath string, configName string) {
	// Настраиваем формат логов в формате JSON
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Загружаем конфигурацию
	cfg, err := config.New(configPath, configName)
	if err != nil {
		logrus.Fatalf("failed to read config: %s", err.Error())
	}

	// Подключаемся к базе данных Postgres
	db, err := repository.NewPostgresDB(&config.Postgres{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Мигрируем базу данных
	if err := Migrate(&cfg.Postgres); err != nil {
		logrus.Fatalf("failed to migrate db: %s", err.Error())
	}

	// Инициализируем пользователя-администратора
	if err := initAdmin(db); err != nil {
		logrus.Fatalf("failed to init admin: %s", err.Error())
	}

	// Создаем репозитории, сервисы и контроллер
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	controller := ctrl.NewController(service)

	// Запускаем сервер в отдельной горутине
	var srv server.Server
	go func() {
		if err := srv.Run(&cfg.Server, controller.Handler.InitRouter()); err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	// Логируем успешный старт приложения
	logrus.Info("todo app started")

	// Ожидаем сигналы для остановки сервера (SIGTERM, SIGINT) Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Логируем завершение работы и закрываем ресурсы
	logrus.Info("shutting down server and database")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
