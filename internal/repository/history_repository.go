// Package repository defines domain-level repository interfaces.
//
// Rules:
//   - No GORM, no *sql.DB, no database-specific types in any interface or
//     method signature.  Callers depend only on Go standard library types
//     and application model types.
//   - Implementations live in internal/database; they are wired together
//     by the application's composition root (cmd/ or runtime/).
package repository

import (
	"context"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/models"
)

// HistoryFilter holds optional filter criteria for composable history queries.
// Zero values are ignored — only non-zero fields are applied.
type HistoryFilter struct {
	MovieID    string
	BatchJobID string
	Operation  string // e.g. "scrape", "organize"
	Status     string // e.g. "success", "failed", "reverted"
	Start      *time.Time
	End        *time.Time
	DryRun     *bool
}

// HistoryRepository is the **domain contract** for history tracking.
//
// Implementations MUST:
//   - Be goroutine-safe (each method is an independent database operation).
//   - Return ErrNotFound (defined below) when a requested record does not exist.
//   - Never expose GORM, *sql.DB, or any driver-specific error types.
//
// The three primary use-cases from the design brief are:
//
//	LogAction   — append a new history entry (write path).
//	GetHistory  — query/filter existing entries (read path).
//	RollbackAction — mark an entry as reverted (mutation path).
type HistoryRepository interface {
	// ── Write ──────────────────────────────────────────────────────────────

	// LogAction persists a new history entry for a file organisation
	// operation.  history.ID is populated by the implementation upon success.
	//
	// Returns an error if history is nil or the write fails.
	LogAction(ctx context.Context, history *models.History) error

	// ── Mutation ───────────────────────────────────────────────────────────

	// RollbackAction marks the entry identified by id as "reverted" and
	// records an optional note in ErrorMessage.
	//
	// Returns ErrNotFound if the id does not exist.
	RollbackAction(ctx context.Context, id uint, note string) error

	// ── Read ───────────────────────────────────────────────────────────────

	// GetHistory returns history entries that match the given filter.
	// limit <= 0 means "return all matching entries".
	// Results are ordered newest-first (created_at DESC).
	GetHistory(ctx context.Context, filter HistoryFilter, limit, offset int) ([]models.History, error)

	// GetByID returns a single history entry by its primary key.
	// Returns ErrNotFound if the id does not exist.
	GetByID(ctx context.Context, id uint) (*models.History, error)

	// GetRecent returns the most recent `limit` entries across all operations.
	GetRecent(ctx context.Context, limit int) ([]models.History, error)

	// ── Aggregation ────────────────────────────────────────────────────────

	// CountByStatus returns the total number of entries with the given status.
	CountByStatus(ctx context.Context, status string) (int64, error)

	// CountByOperation returns the total number of entries for the given operation.
	CountByOperation(ctx context.Context, operation string) (int64, error)

	// TotalCount returns the total number of history entries.
	TotalCount(ctx context.Context) (int64, error)

	// ── Delete ─────────────────────────────────────────────────────────────

	// DeleteByID removes a single entry.  Returns ErrNotFound if it does not exist.
	DeleteByID(ctx context.Context, id uint) error

	// DeleteByMovieID removes all history entries associated with movieID.
	DeleteByMovieID(ctx context.Context, movieID string) error

	// DeleteOlderThan removes all entries created before the given date.
	// This is the primary maintenance operation for retention policies.
	DeleteOlderThan(ctx context.Context, before time.Time) error
}
