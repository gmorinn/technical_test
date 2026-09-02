package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setValidEnv(t *testing.T) {
	t.Helper()

	for key, value := range map[string]string{
		"ENV":               "test",
		"API_DOMAIN":        "localhost",
		"API_TZ":            "Europe/Paris",
		"API_PORT":          "8080",
		"PORT":              "",
		"API_SSL":           "",
		"API_CORS":          "http://localhost:3000",
		"POSTGRES_HOST":     "db",
		"POSTGRES_DB":       "postgres_db",
		"POSTGRES_USER":     "postgres_user",
		"POSTGRES_PASSWORD": "postgres_password",
		"POSTGRES_PORT":     "54320",
	} {
		t.Setenv(key, value)
	}
}

func TestNewConfig_Valid(t *testing.T) {
	setValidEnv(t)

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost:8080", cfg.Host)
	assert.Equal(t, 54320, cfg.Database.Port)
	assert.Equal(t, []string{"http://localhost:3000"}, cfg.Cors)
	assert.Equal(t, "postgres://postgres_user:postgres_password@db:54320/postgres_db?sslmode=disable", cfg.DatabaseURL)
	assert.False(t, cfg.SSL, "API_SSL should default to false when unset")
}

func TestNewConfig_PortOverridesAPIPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("PORT", "9090")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, 9090, cfg.Port)
}

func TestNewConfig_MissingRequiredVariablesAreNamed(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"port", "API_PORT"},
		{"timezone", "API_TZ"},
		{"db host", "POSTGRES_HOST"},
		{"db name", "POSTGRES_DB"},
		{"db user", "POSTGRES_USER"},
		{"db password", "POSTGRES_PASSWORD"},
		{"db port", "POSTGRES_PORT"},
		{"cors", "API_CORS"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(tc.key, "")

			cfg, err := NewConfig()

			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), tc.key, "error must name the offending variable")
		})
	}
}

func TestNewConfig_ReportsEveryProblemAtOnce(t *testing.T) {
	setValidEnv(t)
	t.Setenv("API_PORT", "not-a-number")
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("API_CORS", "")

	_, err := NewConfig()

	require.Error(t, err)
	for _, key := range []string{"API_PORT", "POSTGRES_HOST", "API_CORS"} {
		assert.Contains(t, err.Error(), key)
	}
	assert.Equal(t, 3, strings.Count(err.Error(), "\n"), "one line per problem, after the header")
}

func TestNewConfig_RejectsMalformedPort(t *testing.T) {
	tests := []struct {
		name string
		port string
	}{
		{"not an integer", "eighty-eighty"},
		{"zero", "0"},
		{"negative", "-1"},
		{"above the tcp range", "70000"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv("API_PORT", tc.port)

			_, err := NewConfig()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "API_PORT")
		})
	}
}

func TestNewConfig_RejectsMalformedBool(t *testing.T) {
	setValidEnv(t)
	t.Setenv("API_SSL", "yes-please")

	_, err := NewConfig()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "API_SSL")
}

func TestGetenvSliceString(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []string
	}{
		{"unset", "", nil},
		{"only whitespace", "   ", nil},
		{"only separators", ",,,", []string{}},
		{"single origin", "http://a.test", []string{"http://a.test"}},
		{"multiple origins", "http://a.test,http://b.test", []string{"http://a.test", "http://b.test"}},
		{"padded and trailing comma", " http://a.test , http://b.test , ", []string{"http://a.test", "http://b.test"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TEST_CORS_VALUE", tc.raw)

			assert.Equal(t, tc.expected, getenvSliceString("TEST_CORS_VALUE"))
		})
	}
}
