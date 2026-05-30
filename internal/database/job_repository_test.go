package database

import (
	"testing"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobRepository_Create(t *testing.T) {
	db := newDatabaseTestDB(t)
	repo := NewJobRepository(db)

	job := &models.Job{
		ID:         "test-job-1",
		Status:     string(models.JobStatusRunning),
		TotalFiles: 10,
		Completed:  5,
		Failed:     0,
		Progress:   50.0,
		Files:      `["file1.mp4","file2.mp4"]`,
		StartedAt:  time.Now(),
	}

	err := repo.Create(job)
	require.NoError(t, err)

	found, err := repo.FindByID("test-job-1")
	require.NoError(t, err)
	assert.Equal(t, "test-job-1", found.ID)
	assert.Equal(t, 10, found.TotalFiles)
}

func TestJobRepository_List(t *testing.T) {
	db := newDatabaseTestDB(t)
	repo := NewJobRepository(db)

	now := time.Now()
	jobs := []*models.Job{
		{ID: "job-1", Status: string(models.JobStatusRunning), TotalFiles: 5, Files: "[]", StartedAt: now.Add(-1 * time.Hour)},
		{ID: "job-2", Status: string(models.JobStatusCompleted), TotalFiles: 3, Files: "[]", StartedAt: now},
	}

	for _, j := range jobs {
		require.NoError(t, repo.Create(j))
	}

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Equal(t, "job-2", list[0].ID)
}

func TestJobRepository_Delete(t *testing.T) {
	db := newDatabaseTestDB(t)
	repo := NewJobRepository(db)

	job := &models.Job{
		ID:         "to-delete",
		Status:     string(models.JobStatusCompleted),
		TotalFiles: 1,
		Files:      "[]",
		StartedAt:  time.Now(),
	}
	require.NoError(t, repo.Create(job))

	err := repo.Delete("to-delete")
	require.NoError(t, err)

	_, err = repo.FindByID("to-delete")
	assert.Error(t, err)
}

func TestJobRepository_DeleteOrganizedOlderThan(t *testing.T) {
	db := newDatabaseTestDB(t)
	repo := NewJobRepository(db)

	now := time.Now()
	twoDaysAgo := now.Add(-48 * time.Hour)

	organizedOld := &models.Job{
		ID:          "organized-old",
		Status:      string(models.JobStatusOrganized),
		TotalFiles:  1,
		Files:       "[]",
		StartedAt:   twoDaysAgo.Add(-1 * time.Hour),
		OrganizedAt: &twoDaysAgo,
	}
	organizedRecent := &models.Job{
		ID:          "organized-recent",
		Status:      string(models.JobStatusOrganized),
		TotalFiles:  1,
		Files:       "[]",
		StartedAt:   now.Add(-1 * time.Hour),
		OrganizedAt: ptrTime(now.Add(-12 * time.Hour)),
	}

	require.NoError(t, repo.Create(organizedOld))
	require.NoError(t, repo.Create(organizedRecent))

	err := repo.DeleteOrganizedOlderThan(now.Add(-24 * time.Hour))
	require.NoError(t, err)

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "organized-recent", list[0].ID)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
