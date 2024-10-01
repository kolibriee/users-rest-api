package config

import (
	"errors"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Основная структура конфигурации
type Config struct {
	Postgres Postgres // Конфигурация PostgreSQL
	Server   Server   `mapstructure:"server"` // Конфигурация сервера
}

// Структура конфигурации сервера
type Server struct {
	Port           string        `mapstructure:"port"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
}

// Структура конфигурации PostgreSQL
type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Функция для создания новой конфигурации
func New(path string, fileName string) (*Config, error) {
	var cfg *Config

	// Загружаем переменные окружения из файла .env
	godotenv.Load()

	// Настраиваем viper для чтения конфигурационного файла
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)

	// Читаем конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	// Разворачиваем конфигурацию в структуру
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.New("failed to unmarshal config: " + err.Error())
	}

	// Обрабатываем переменные окружения для конфигурации базы данных
	if err := envconfig.Process("DB", &cfg.Postgres); err != nil {
		return nil, errors.New("failed to process env variables: " + err.Error())
	}

	return cfg, nil
}
