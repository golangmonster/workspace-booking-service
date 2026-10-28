//go:build integration

package testdb

import (
	"context"
	"crypto/rand"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Postgres struct {
	admin      *pgxpool.Pool
	migrations fs.FS
}

type Database struct {
	Pool *pgxpool.Pool
	db   *sql.DB
}

// Start resolves migrationsPath relative to the project root containing go.mod.
func Start(t *testing.T, migrationsPath string) *Postgres {
	t.Helper()

	root, err := os.Getwd()
	require.NoError(t, err)

	for {
		_, err := os.Stat(filepath.Join(root, "go.mod"))
		if err == nil {
			break
		}
		require.ErrorIs(t, err, os.ErrNotExist)

		parent := filepath.Dir(root)
		require.NotEqual(t, root, parent, "project root with go.mod not found")
		root = parent
	}

	dir := filepath.Join(root, migrationsPath)
	require.DirExists(t, dir)

	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(time.Minute)),
	)
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return &Postgres{admin: pool, migrations: os.DirFS(dir)}
}

// NewDatabase creates an isolated database, applies migrations and loads fixtures.
func (p *Postgres) NewDatabase(t *testing.T, initPath string) *Database {
	t.Helper()

	name := "test_" + rand.Text()
	identifier := pgx.Identifier{name}.Sanitize()

	_, err := p.admin.Exec(t.Context(), "CREATE DATABASE "+identifier)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := p.admin.Exec(ctx, "DROP DATABASE "+identifier+" WITH (FORCE)")
		if err != nil {
			t.Errorf("drop test database: %v", err)
		}
	})

	cfg := p.admin.Config().Copy()
	cfg.ConnConfig.Database = name

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	db := stdlib.OpenDB(*cfg.ConnConfig)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	provider, err := goose.NewProvider(goose.DialectPostgres, db, p.migrations)
	require.NoError(t, err)

	_, err = provider.Up(t.Context())
	require.NoError(t, err)

	database := &Database{Pool: pool, db: db}
	database.load(t, initPath)

	return database
}
