package main

import "github.com/kolibriee/users-rest-api/internal/app"

const (
	configName = "config"
	configsDir = "configs"
)

func main() {
	app.Run(configsDir, configName)
}
