// Package testdb prepares isolated PostgreSQL databases for integration tests.
package testdb

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golangmonster/workspace-booking-service/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// StartPostgres starts one container for a parent test and its subtests.
func StartPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("postgres"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err, "start postgres")

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "get postgres DSN")

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "open postgres pool")
	t.Cleanup(pool.Close)

	require.NoError(t, pool.Ping(ctx), "connect to postgres")

	return pool
}

// NewDatabase gives each test its own database, connection pool, and migrations.
func NewDatabase(t *testing.T, postgres *pgxpool.Pool) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	name := "test_" + rand.Text()
	identifier := pgx.Identifier{name}.Sanitize()

	_, err := postgres.Exec(ctx, "CREATE DATABASE "+identifier)
	require.NoError(t, err, "create test database")

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := postgres.Exec(cleanupCtx, "DROP DATABASE "+identifier+" WITH (FORCE)")
		if err != nil {
			t.Errorf("drop test database: %v", err)
		}
	})

	cfg := postgres.Config()
	cfg.ConnConfig.Database = name

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err, "open test database pool")
	t.Cleanup(pool.Close)

	db := stdlib.OpenDB(*cfg.ConnConfig)
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.Files)
	require.NoError(t, err, "create migration provider")

	_, err = provider.Up(ctx)
	require.NoError(t, err, "apply migrations")

	return pool
}
