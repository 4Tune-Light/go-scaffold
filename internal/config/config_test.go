package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config not found")
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("")
	assert.Error(t, err)
}

func TestPostgresConfig_GeneratesDSN(t *testing.T) {
	cfg := &Config{
		Database: Database{
			Postgres: Postgres{
				Host: "myhost", Port: 5432, User: "myuser",
				Password: "mypass", DBName: "mydb", SSLMode: "require",
			},
		},
	}

	pg := cfg.PostgresConfig()
	assert.Contains(t, pg.DSN, "myhost:5432")
	assert.Contains(t, pg.DSN, "myuser")
	assert.Contains(t, pg.DSN, "mydb")
	assert.Contains(t, pg.DSN, "sslmode=require")
}

func TestRedisConfig_GeneratesAddr(t *testing.T) {
	cfg := &Config{
		Redis: Redis{Host: "myhost", Port: 6379, DB: 1},
	}

	rc := cfg.RedisConfig()
	assert.Equal(t, "myhost:6379", rc.Addr)
	assert.Equal(t, 1, rc.DB)
}

func TestValidate_Invalid(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "app.name is required")
	assert.Contains(t, err.Error(), "server.http.port is required")
}

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		App: App{Name: "test"},
		Server: Server{HTTP: ServerHTTP{Port: 8080}},
		Database: Database{Postgres: Postgres{Host: "localhost", User: "u", DBName: "d"}},
	}
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestLoad_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	yaml := `app:
  name: testapp
server:
  http:
    port: 8080
database:
  postgres:
    host: localhost
    user: test
    dbname: testdb
    sslmode: disable
redis:
  host: localhost
  port: 6379`

	yamlPath := filepath.Join(dir, "config.yaml")
	writeFile(t, yamlPath, yaml)

	cfg, err := Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "testapp", cfg.App.Name)
	assert.Equal(t, 8080, cfg.Server.HTTP.Port)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
