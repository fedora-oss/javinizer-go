package database

import (
	"fmt"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/models"
)

type HistoryRepository struct {
	*BaseRepository[models.History, uint]
}

func NewHistoryRepository(db *DB) *HistoryRepository {
	return &HistoryRepository{
		BaseRepository: NewBaseRepository[models.History, uint](
			db, "history",
			func(h models.History) string { return fmt.Sprintf("%d", h.ID) },
			WithDefaultOrder[models.History, uint]("created_at DESC"),
			WithNewEntity[models.History, uint](func() models.History { return models.History{} }),
		),
	}
}

func (r *HistoryRepository) Create(history *models.History) error {
	return r.BaseRepository.Create(history)
}

func (r *HistoryRepository) FindByID(id uint) (*models.History, error) {
	return r.BaseRepository.FindByID(id)
}

func (r *HistoryRepository) FindByMovieID(movieID string) ([]models.History, error) {
	var history []models.History
	err := r.GetDB().Where("movie_id = ?", movieID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("history for movie %s", movieID), err)
	}
	return history, nil
}

func (r *HistoryRepository) FindByOperation(operation string, limit int) ([]models.History, error) {
	var history []models.History
	query := r.GetDB().Where("operation = ?", operation).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("history by operation %s", operation), err)
	}
	return history, nil
}

func (r *HistoryRepository) FindByStatus(status string, limit int) ([]models.History, error) {
	var history []models.History
	query := r.GetDB().Where("status = ?", status).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("history by status %s", status), err)
	}
	return history, nil
}

func (r *HistoryRepository) FindRecent(limit int) ([]models.History, error) {
	var history []models.History
	err := r.GetDB().Order("created_at DESC").Limit(limit).Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", "recent history", err)
	}
	return history, nil
}

func (r *HistoryRepository) FindByDateRange(start, end time.Time) ([]models.History, error) {
	var history []models.History
	// Use the dialect helper so the query works on SQLite, PostgreSQL, and MySQL.
	clause, args := r.GetDB().Dialect.BetweenDateTimeExpr("created_at", start, end)
	err := r.GetDB().Where(clause, args...).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", "history by date range", err)
	}
	return history, nil
}

func (r *HistoryRepository) Count() (int64, error) {
	return r.BaseRepository.Count()
}

func (r *HistoryRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.GetDB().Model(&models.History{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, wrapDBErr("count", fmt.Sprintf("history by status %s", status), err)
	}
	return count, nil
}

func (r *HistoryRepository) CountByOperation(operation string) (int64, error) {
	var count int64
	err := r.GetDB().Model(&models.History{}).Where("operation = ?", operation).Count(&count).Error
	if err != nil {
		return 0, wrapDBErr("count", fmt.Sprintf("history by operation %s", operation), err)
	}
	return count, nil
}

func (r *HistoryRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(id)
}

func (r *HistoryRepository) DeleteByMovieID(movieID string) error {
	if err := r.GetDB().Where("movie_id = ?", movieID).Delete(&models.History{}).Error; err != nil {
		return wrapDBErr("delete", fmt.Sprintf("history for movie %s", movieID), err)
	}
	return nil
}

func (r *HistoryRepository) DeleteOlderThan(date time.Time) error {
	// Use the dialect helper so the query works on SQLite, PostgreSQL, and MySQL.
	clause, args := r.GetDB().Dialect.BeforeDateTimeExpr("created_at", date)
	if err := r.GetDB().Where(clause, args...).Delete(&models.History{}).Error; err != nil {
		return wrapDBErr("delete", "history older than date", err)
	}
	return nil
}

func (r *HistoryRepository) List(limit, offset int) ([]models.History, error) {
	return r.BaseRepository.List(limit, offset)
}

func (r *HistoryRepository) FindByBatchJobID(batchJobID string) ([]models.History, error) {
	var history []models.History
	err := r.GetDB().Where("batch_job_id = ?", batchJobID).Order("created_at ASC").Find(&history).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("history for batch job %s", batchJobID), err)
	}
	return history, nil
}
