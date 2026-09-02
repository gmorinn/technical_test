package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kataras/golog"
)

type API struct {
	Mode   string
	Domain string
	TZ     string
	SSL    bool
	Host   string
	Cors   []string
	Port   int
	// Cert        string
	// Key         string
	DatabaseURL string
	Database    Database
	URLFront    string
}

type Database struct {
	Host     string
	Database string
	User     string
	Password string
	Port     int
	SSLMode  string
}

func NewConfig() (*API, error) {
	var (
		config   API
		problems []error
	)

	add := func(format string, args ...any) {
		problems = append(problems, fmt.Errorf(format, args...))
	}

	requireString := func(key string) string {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			add("%s is required but not set", key)
		}
		return value
	}

	requirePort := func(key string) int {
		raw := strings.TrimSpace(os.Getenv(key))
		if raw == "" {
			add("%s is required but not set", key)
			return 0
		}
		value, err := strconv.Atoi(raw)
		if err != nil {
			add("%s must be an integer, got %q", key, raw)
			return 0
		}
		if value < 1 || value > 65535 {
			add("%s must be a TCP port between 1 and 65535, got %d", key, value)
			return 0
		}
		return value
	}

	optionalBool := func(key string, fallback bool) bool {
		raw := strings.TrimSpace(os.Getenv(key))
		if raw == "" {
			return fallback
		}
		value, err := strconv.ParseBool(raw)
		if err != nil {
			add("%s must be a boolean, got %q", key, raw)
			return fallback
		}
		return value
	}

	config.Mode = os.Getenv("ENV")
	config.Domain = os.Getenv("API_DOMAIN")
	config.URLFront = os.Getenv("URL_FRONT")

	// TZ is interpolated into the Postgres DSN; an empty value makes the server
	// reject the connection with an opaque parameter error.
	config.TZ = requireString("API_TZ")

	portKey := "API_PORT"
	if strings.TrimSpace(os.Getenv("PORT")) != "" {
		portKey = "PORT"
	}
	config.Port = requirePort(portKey)

	config.SSL = optionalBool("API_SSL", false)

	if config.Domain != "localhost" {
		config.Host = config.Domain
	} else {
		config.Host = fmt.Sprintf("%s:%v", config.Domain, config.Port)
	}

	config.Database.Host = requireString("POSTGRES_HOST")
	config.Database.Database = requireString("POSTGRES_DB")
	config.Database.User = requireString("POSTGRES_USER")
	config.Database.Password = requireString("POSTGRES_PASSWORD")
	config.Database.Port = requirePort("POSTGRES_PORT")

	config.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Database)

	config.Cors = getenvSliceString("API_CORS")
	if len(config.Cors) == 0 {
		add("API_CORS is required but not set (comma-separated list of allowed origins)")
	}

	if len(problems) > 0 {
		return nil, fmt.Errorf("invalid configuration:\n%w", errors.Join(problems...))
	}

	if config.Mode == "dev" {
		golog.SetLevel("debug")
	}

	return &config, nil
}

// getenvSliceString reads a comma-separated variable, trimming blanks. It
// returns nil rather than []string{""} when the variable is unset or empty, so
// callers can use len() to detect absence.
func getenvSliceString(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			values = append(values, part)
		}
	}
	return values
}
