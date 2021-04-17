package main

import (
	"github.com/hrshadhin/fiber-go-boilerplate/cmd/server"
	_ "github.com/hrshadhin/fiber-go-boilerplate/docs" // load API Docs files (Swagger)
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
)

// @title Fiber Go API
// @version 1.0
// @description Fiber go web framework based REST API boilerplate
// @contact.name H.R. Shadhin
// @contact.email dev@hrshadhin.me
// @termsOfService
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:5000
// @BasePath /api
func main() {

	// setup various configuration for app
	config.LoadAllConfigs(".env")

	server.Serve()
}
