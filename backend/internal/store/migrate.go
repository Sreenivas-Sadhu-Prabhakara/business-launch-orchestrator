package store

import (
	"context"
	"fmt"
	"io/fs"
	"sort"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/migrations"
)

// Migrate applies every embedded .sql migration in lexical order. Migrations
// are written to be idempotent (CREATE TABLE IF NOT EXISTS ...), so running
// them on every startup is safe.
func (s *Store) Migrate(ctx context.Context) error {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 4 && e.Name()[len(e.Name())-4:] == ".sql" {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		sqlBytes, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if _, err := s.pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}
