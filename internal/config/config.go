package config

import "os"

// AppConfig is a struct that contains the env vars for the application.
type AppConfig struct {
	// Environment in which the app is running.
	Env string `env:"APP_ENV" default:"dev"`
	// Log format (json | text).
	LogFormat string `env:"LOG_FORMAT" default:"json"`
	// Log level (debug | info | warn | error).
	LogLevel string `env:"LOG_LEVEL" default:"info"`
	// HTTP port the server listens on.
	Port int `env:"PORT"`
	// Public key for JWT signing.
	AuthPublKey string `env:"AUTH_PUBLIC_KEY"`
	// Private key for JWT signing.
	AuthPrivKey string `env:"AUTH_PRIVATE_KEY"`
	// Connection string to the Postgres database.
	PostgresConnStr string `env:"POSTGRES_CONN_STR"`
}

var appConfig AppConfig

func init() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	err := Parse(&appConfig, envFile)
	if err != nil {
		panic(err)
	}
}

// Get returns the AppConfig, which is populated at `init` and panics if there is an error parsing env vars.
func Get() *AppConfig {
	return &appConfig
}
