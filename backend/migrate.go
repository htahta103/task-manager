package main

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func runMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := migrationsFS.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	type mig struct {
		name string
		path string
	}
	var m []mig
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		m = append(m, mig{name: name, path: migrationsDir + "/" + name})
	}
	sort.Slice(m, func(i, j int) bool { return m[i].name < m[j].name })

	for _, mi := range m {
		b, err := migrationsFS.ReadFile(mi.path)
		if err != nil {
			return err
		}
		sql := strings.TrimSpace(string(b))
		if sql == "" {
			continue
		}
		if _, err := pool.Exec(ctx, sql); err != nil {
			return fmt.Errorf("migration %s failed: %w", mi.name, err)
		}
	}
	return nil
}
