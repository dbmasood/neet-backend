package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		GRPC    GRPC
		RMQ     RMQ
		NATS    NATS
		Metrics Metrics
		Swagger Swagger
		JWT     JWT
		Admin   Admin
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// GRPC -.
	GRPC struct {
		Port string `env:"GRPC_PORT,required"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// NATS -.
	NATS struct {
		ServerExchange string `env:"NATS_RPC_SERVER,required"`
		URL            string `env:"NATS_URL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
	JWT struct {
		UserSecret      string `env:"JWT_USER_SECRET,required"`
		AdminSecret     string `env:"JWT_ADMIN_SECRET,required"`
		TokenTTLMinutes int    `env:"JWT_TOKEN_TTL_MINUTES" envDefault:"1440"`
	}

	Admin struct {
		Username     string   `env:"ADMIN_USERNAME,required"`
		Password     string   `env:"ADMIN_PASSWORD,required"`
		DisplayName  string   `env:"ADMIN_DISPLAY_NAME" envDefault:"Super Admin"`
		Email        string   `env:"ADMIN_EMAIL,required"`
		UserID       string   `env:"ADMIN_USER_ID" envDefault:"00000000-0000-0000-0000-000000000001"`
		PrimaryExam  string   `env:"ADMIN_PRIMARY_EXAM" envDefault:"NEET_PG"`
		Role         string   `env:"ADMIN_ROLE" envDefault:"ADMIN"`
		Permissions  []string `env:"ADMIN_PERMISSIONS" envDefault:"subjects.read,subjects.write" envSeparator:","`
		CreatedAtISO string   `env:"ADMIN_CREATED_AT"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
