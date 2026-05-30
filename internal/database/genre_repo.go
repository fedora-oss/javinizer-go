package database

import (
	"fmt"

	"github.com/fedora-oss/javinizer-go/internal/models"
)

type GenreRepository struct {
	*BaseRepository[models.Genre, uint]
}

func NewGenreRepository(db *DB) *GenreRepository {
	return &GenreRepository{
		BaseRepository: NewBaseRepository[models.Genre, uint](
			db, "genre",
			func(g models.Genre) string { return g.Name },
			WithNewEntity[models.Genre, uint](func() models.Genre { return models.Genre{} }),
		),
	}
}

func (r *GenreRepository) FindOrCreate(name string) (*models.Genre, error) {
	var genre models.Genre
	err := r.GetDB().FirstOrCreate(&genre, models.Genre{Name: name}).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("genre %s", name), err)
	}
	return &genre, nil
}

func (r *GenreRepository) List() ([]models.Genre, error) {
	return r.ListAll()
}
