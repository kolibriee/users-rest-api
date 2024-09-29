package app

import (
	"context"
	"fmt"
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
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg, err := config.New(configPath, configName)
	if err != nil {
		logrus.Fatalf("failed to read config: %s", err.Error())
	}
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
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	controller := ctrl.NewController(service)
	var srv server.Server
	go func() {
		if err := srv.Run(&cfg.Server, controller.Handler.InitRouter()); err != nil {
			logrus.Fatalf("error occured while runnirest server: %s", err.Error())
		}
	}()
	fmt.Println(os.Getenv("TOKEN_SECRET_KEY"))
	logrus.Info("todo app started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Info("shutting down server and database")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
