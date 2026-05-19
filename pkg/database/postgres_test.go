package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration_NewPostgresPool(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "postgres:16-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	cfg := PostgresConfig{
		DSN: "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable",
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	}

	pool, err := NewPostgresPool(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	err = pool.Ping(ctx)
	assert.NoError(t, err)
}

func TestIntegration_NewRedisClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	host, _ := redisContainer.Host(ctx)
	port, _ := redisContainer.MappedPort(ctx, "6379")

	rdb, err := NewRedisClient(ctx, RedisConfig{
		Addr: host + ":" + port.Port(),
	})
	require.NoError(t, err)
	defer rdb.Close()

	_, err = rdb.Ping(ctx).Result()
	assert.NoError(t, err)
}


