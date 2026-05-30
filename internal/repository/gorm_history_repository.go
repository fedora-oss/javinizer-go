// Package repository contains the GORM-backed implementation of the
// HistoryRepository domain interface.
//
// This file is the ONLY place in the codebase where GORM internals are
// permitted to appear in history-related code.  All callers depend only on
// the HistoryRepository interface defined in history_repository.go.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/database"
	"github.com/fedora-oss/javinizer-go/internal/models"
	"gorm.io/gorm"
)

// Compile-time assertion: gormHistoryRepo must satisfy HistoryRepository.
var _ HistoryRepository = (*gormHistoryRepo)(nil)

// gormHistoryRepo is the GORM-backed implementation of HistoryRepository.
// It is unexported so that the rest of the application is forced to depend on
// the interface, not the concrete type.
type gormHistoryRepo struct {
	// db is the shared database wrapper that carries both the *gorm.DB
	// connection and the SQL Dialect for the active backend.
	db *database.DB
}

// NewGORMHistoryRepository constructs a gormHistoryRepo and returns it as the
// HistoryRepository interface.  The factory hides the concrete type from
// callers — this is the Factory Pattern requested in the brief.
//
// Usage (in your composition root / runtime wiring):
//
//	db, err := database.New(cfg)
//	var histRepo repository.HistoryRepository = repository.NewGORMHistoryRepository(db)
func NewGORMHistoryRepository(db *database.DB) HistoryRepository {
	return &gormHistoryRepo{db: db}
}

// ── Write ─────────────────────────────────────────────────────────────────────

// LogAction persists a new history entry.
// history.ID is set by GORM after a successful INSERT.
func (r *gormHistoryRepo) LogAction(ctx context.Context, history *models.History) error {
	if history == nil {
		return fmt.Errorf("LogAction: history must not be nil")
	}
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		return fmt.Errorf("LogAction: %w", err)
	}
	return nil
}

// ── Mutation ──────────────────────────────────────────────────────────────────

// RollbackAction marks the entry identified by id as "reverted" and stores
// the note in ErrorMessage.  Returns ErrNotFound when the id is missing.
func (r *gormHistoryRepo) RollbackAction(ctx context.Context, id uint, note string) error {
	result := r.db.WithContext(ctx).
		Model(&models.History{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "reverted",
			"error_message": note,
		})

	if result.Error != nil {
		return fmt.Errorf("RollbackAction %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		// No row was updated — the ID does not exist.
		return fmt.Errorf("RollbackAction %d: %w", id, ErrNotFound)
	}
	return nil
}

// ── Read ──────────────────────────────────────────────────────────────────────

// GetHistory returns history entries that match the given filter, ordered
// newest-first.  All filter fields are optional; zero values are skipped.
func (r *gormHistoryRepo) GetHistory(ctx context.Context, filter HistoryFilter, limit, offset int) ([]models.History, error) {
	q := r.db.WithContext(ctx).Order("created_at DESC")

	// Apply each filter field only when it carries a non-zero value.
	if filter.MovieID != "" {
		q = q.Where("movie_id = ?", filter.MovieID)
	}
	if filter.BatchJobID != "" {
		q = q.Where("batch_job_id = ?", filter.BatchJobID)
	}
	if filter.Operation != "" {
		q = q.Where("operation = ?", filter.Operation)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.DryRun != nil {
		q = q.Where("dry_run = ?", *filter.DryRun)
	}

	// Date-range filtering uses the dialect helper so the generated SQL is
	// correct for the active backend (SQLite datetime() vs. native timestamps).
	if filter.Start != nil && filter.End != nil {
		clause, args := r.db.Dialect.BetweenDateTimeExpr("created_at", *filter.Start, *filter.End)
		q = q.Where(clause, args...)
	} else if filter.Start != nil {
		// Only a lower bound — use a simple comparison.
		q = q.Where("created_at >= ?", filter.Start.UTC())
	} else if filter.End != nil {
		// Only an upper bound.
		clause, args := r.db.Dialect.BeforeDateTimeExpr("created_at", *filter.End)
		q = q.Where(clause, args...)
	}

	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}

	var results []models.History
	if err := q.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("GetHistory: %w", err)
	}
	return results, nil
}

// GetByID returns a single history entry by its primary key.
// Returns ErrNotFound when the record does not exist.
func (r *gormHistoryRepo) GetByID(ctx context.Context, id uint) (*models.History, error) {
	var h models.History
	err := r.db.WithContext(ctx).First(&h, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("GetByID %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("GetByID %d: %w", id, err)
	}
	return &h, nil
}

// GetRecent returns the `limit` most recent history entries across all
// operations, ordered newest-first.
func (r *gormHistoryRepo) GetRecent(ctx context.Context, limit int) ([]models.History, error) {
	var results []models.History
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("GetRecent: %w", err)
	}
	return results, nil
}

// ── Aggregation ───────────────────────────────────────────────────────────────

// CountByStatus returns the number of history entries with the given status.
func (r *gormHistoryRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.History{}).
		Where("status = ?", status).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("CountByStatus %q: %w", status, err)
	}
	return count, nil
}

// CountByOperation returns the number of history entries for the given operation.
func (r *gormHistoryRepo) CountByOperation(ctx context.Context, operation string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.History{}).
		Where("operation = ?", operation).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("CountByOperation %q: %w", operation, err)
	}
	return count, nil
}

// TotalCount returns the total number of history entries.
func (r *gormHistoryRepo) TotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.History{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("TotalCount: %w", err)
	}
	return count, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

// DeleteByID removes a single history entry.
// Returns ErrNotFound when the id does not exist.
func (r *gormHistoryRepo) DeleteByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.History{}, id)
	if result.Error != nil {
		return fmt.Errorf("DeleteByID %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("DeleteByID %d: %w", id, ErrNotFound)
	}
	return nil
}

// DeleteByMovieID removes all history entries for the given movie ID.
func (r *gormHistoryRepo) DeleteByMovieID(ctx context.Context, movieID string) error {
	if err := r.db.WithContext(ctx).Where("movie_id = ?", movieID).Delete(&models.History{}).Error; err != nil {
		return fmt.Errorf("DeleteByMovieID %q: %w", movieID, err)
	}
	return nil
}

// DeleteOlderThan removes all history entries created before `before`.
// The dialect helper produces engine-correct SQL for the comparison.
func (r *gormHistoryRepo) DeleteOlderThan(ctx context.Context, before time.Time) error {
	clause, args := r.db.Dialect.BeforeDateTimeExpr("created_at", before)
	if err := r.db.WithContext(ctx).Where(clause, args...).Delete(&models.History{}).Error; err != nil {
		return fmt.Errorf("DeleteOlderThan: %w", err)
	}
	return nil
}
