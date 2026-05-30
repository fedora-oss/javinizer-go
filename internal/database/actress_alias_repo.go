package database

import (
	"fmt"

	"github.com/fedora-oss/javinizer-go/internal/models"
)

type ActressAliasRepository struct {
	*BaseRepository[models.ActressAlias, uint]
}

func NewActressAliasRepository(db *DB) *ActressAliasRepository {
	return &ActressAliasRepository{
		BaseRepository: NewBaseRepository[models.ActressAlias, uint](
			db, "actress alias",
			func(a models.ActressAlias) string { return a.AliasName },
			WithNewEntity[models.ActressAlias, uint](func() models.ActressAlias { return models.ActressAlias{} }),
		),
	}
}

func (r *ActressAliasRepository) Create(alias *models.ActressAlias) error {
	return r.BaseRepository.Create(alias)
}

func (r *ActressAliasRepository) Upsert(alias *models.ActressAlias) error {
	existing, err := r.FindByAliasName(alias.AliasName)
	if err != nil {
		if !isRecordNotFound(err) {
			return err
		}
		return r.Create(alias)
	}

	alias.ID = existing.ID
	alias.CreatedAt = existing.CreatedAt
	if err := r.GetDB().Save(alias).Error; err != nil {
		return wrapDBErr("update", fmt.Sprintf("actress alias %s", alias.AliasName), err)
	}
	return nil
}

func (r *ActressAliasRepository) FindByAliasName(aliasName string) (*models.ActressAlias, error) {
	var alias models.ActressAlias
	err := r.GetDB().First(&alias, "alias_name = ?", aliasName).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("actress alias %s", aliasName), err)
	}
	return &alias, nil
}

func (r *ActressAliasRepository) FindByCanonicalName(canonicalName string) ([]models.ActressAlias, error) {
	var aliases []models.ActressAlias
	err := r.GetDB().Where("canonical_name = ?", canonicalName).Find(&aliases).Error
	if err != nil {
		return nil, wrapDBErr("find", fmt.Sprintf("actress aliases for %s", canonicalName), err)
	}
	return aliases, nil
}

func (r *ActressAliasRepository) List() ([]models.ActressAlias, error) {
	return r.ListAll()
}

func (r *ActressAliasRepository) Delete(aliasName string) error {
	if err := r.GetDB().Delete(&models.ActressAlias{}, "alias_name = ?", aliasName).Error; err != nil {
		return wrapDBErr("delete", fmt.Sprintf("actress alias %s", aliasName), err)
	}
	return nil
}

func (r *ActressAliasRepository) GetAliasMap() (map[string]string, error) {
	aliases, err := r.List()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, a := range aliases {
		result[a.AliasName] = a.CanonicalName
	}
	return result, nil
}
