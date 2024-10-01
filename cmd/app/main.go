package main

import "github.com/kolibriee/users-rest-api/internal/app"

const (
	configName = "config"
	configsDir = "configs"
)

//	@title			Users REST API
//	@version		1.0
//	@description	API Server for Users Service for your app

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	app.Run(configsDir, configName)
}
