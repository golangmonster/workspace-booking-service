package testdb

import (
	"net"
	"net/url"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver used by dsunit.
	"github.com/stretchr/testify/require"
	"github.com/viant/dsc"
	"github.com/viant/dsunit"
)

// NewFixtures connects dsunit to the test database under its database name.
// Each test owns its service and connections; no global dsunit registry is used.
func NewFixtures(t *testing.T, pool *pgxpool.Pool) dsunit.Service {
	t.Helper()

	cfg := pool.Config().ConnConfig
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))),
		Path:     "/" + cfg.Database,
		RawQuery: "sslmode=disable",
	}

	fixtures := dsunit.New()
	registered := fixtures.Register(dsunit.NewRegisterRequest(cfg.Database, &dsc.Config{
		DriverName: "postgres",
		Descriptor: dsn.String(),
	}))
	require.NoError(t, registered.Error(), "register fixture database")

	t.Cleanup(func() {
		manager := fixtures.Registry().Get(cfg.Database)
		if err := manager.ConnectionProvider().Close(); err != nil {
			t.Errorf("close fixture connections: %v", err)
		}
	})

	return fixtures
}
