package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	pool *pgxpool.Pool
}

func NewMigrator(pool *pgxpool.Pool) *Migrator {
	return &Migrator{pool: pool}
}

func (m *Migrator) Run(ctx context.Context, migrations []string) error {
	_, err := m.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS sivc_migrations (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) UNIQUE NOT NULL,
		executed_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("table migrations: %w", err)
	}

	for i, migration := range migrations {
		var exists bool
		m.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM sivc_migrations WHERE filename = $1)`, fmt.Sprintf("migration_%d", i)).Scan(&exists)
		if exists {
			continue
		}
		if _, err := m.pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration %d: %w", i, err)
		}
		m.pool.Exec(ctx, `INSERT INTO sivc_migrations (filename) VALUES ($1)`, fmt.Sprintf("migration_%d", i))
	}
	return nil
}
