package config

import (
	"os"
	"strconv"
	"time"
)

// App holds the App configuration
type App struct {
	Host        string
	Port        int
	Debug       bool
	ReadTimeout time.Duration

	// JWT Conf
	JWTSecretKey                string
	JWTSecretExpireMinutesCount int
}

var app = &App{}

// AppCfg returns the default App configuration
func AppCfg() *App {
	return app
}

// LoadApp loads App configuration
func LoadApp() {
	app.Host = os.Getenv("APP_HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	app.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	timeOut, _ := strconv.Atoi(os.Getenv("APP_READ_TIMEOUT"))
	app.ReadTimeout = time.Duration(timeOut) * time.Second

	app.JWTSecretKey = os.Getenv("APP_PORT")
	app.JWTSecretExpireMinutesCount, _ = strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

}
