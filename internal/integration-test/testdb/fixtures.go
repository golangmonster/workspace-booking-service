//go:build integration

package testdb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func (d *Database) load(t *testing.T, path string) {
	t.Helper()

	loader, err := testfixtures.New(
		testfixtures.Database(d.db),
		testfixtures.Dialect("postgres"),
		testfixtures.FilesMultiTables(path),
		testfixtures.Location(time.UTC),
	)
	require.NoError(t, err)
	require.NoError(t, loader.Load(), "load %s", path)
}

// Assert checks all tables in the current schema, ignoring row order.
// Tables omitted from path must be empty; the Goose version table is excluded.
func (d *Database) Assert(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var expected map[string][]map[string]any
	require.NoError(t, yaml.Unmarshal(data, &expected))
	require.NotNil(t, expected, "use {} to expect an empty database")

	result, err := d.Pool.Query(t.Context(), `
		SELECT tablename FROM pg_tables
		WHERE schemaname = current_schema() AND tablename <> $1
		ORDER BY tablename
	`, goose.DefaultTablename)
	require.NoError(t, err)

	tables, err := pgx.CollectRows(result, pgx.RowTo[string])
	require.NoError(t, err)

	for table := range expected {
		require.True(t, slices.Contains(tables, table), "unknown table in assert file: %s", table)
	}

	dir := t.TempDir()
	dumper, err := testfixtures.NewDumper(
		testfixtures.DumpDatabase(d.db),
		testfixtures.DumpDialect("postgres"),
		testfixtures.DumpDirectory(dir),
		testfixtures.DumpTables(tables...),
	)
	require.NoError(t, err)
	require.NoError(t, dumper.Dump())

	for _, table := range tables {
		data, err := os.ReadFile(filepath.Join(dir, table+".yml"))
		require.NoError(t, err)

		var actual []map[string]any
		require.NoError(t, yaml.Unmarshal(data, &actual))

		// Dynamic timestamps are optional; explicitly listed timestamps are checked.
		for _, column := range []string{"created_at", "updated_at"} {
			if slices.ContainsFunc(expected[table], func(row map[string]any) bool {
				_, exists := row[column]

				return exists
			}) {
				continue
			}

			for _, row := range actual {
				delete(row, column)
			}
		}

		normalizeJSON(expected[table])
		normalizeJSON(actual)
		require.ElementsMatch(t, expected[table], actual, "table %s; actual dump: %s", table, dir)
	}
}

// JSON stored in text columns is compared by content rather than formatting.
func normalizeJSON(rows []map[string]any) {
	for _, row := range rows {
		for column, value := range row {
			text, ok := value.(string)
			if !ok {
				continue
			}

			text = strings.TrimSpace(text)
			if !strings.HasPrefix(text, "{") && !strings.HasPrefix(text, "[") {
				continue
			}

			var decoded any
			if json.Unmarshal([]byte(text), &decoded) == nil {
				row[column] = decoded
			}
		}
	}
}
