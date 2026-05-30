package token

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/fedora-oss/javinizer-go/internal/database"
	"github.com/fedora-oss/javinizer-go/internal/logging"
	"github.com/fedora-oss/javinizer-go/internal/models"
)

type TokenService struct {
	repo database.ApiTokenRepositoryInterface
}

func NewTokenService(repo database.ApiTokenRepositoryInterface) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) Create(name string) (*models.ApiToken, string, error) {
	fullToken, prefix, err := GenerateToken()
	if err != nil {
		return nil, "", err
	}

	hash := HashToken(fullToken)

	id, err := generateUUID()
	if err != nil {
		return nil, "", err
	}

	token := &models.ApiToken{
		ID:          id,
		Name:        name,
		TokenHash:   hash,
		TokenPrefix: prefix,
	}

	if err := s.repo.Create(token); err != nil {
		return nil, "", err
	}

	return token, fullToken, nil
}

func (s *TokenService) Revoke(id string) error {
	return s.repo.Revoke(id)
}

func (s *TokenService) List() ([]models.ApiToken, error) {
	return s.repo.ListActive()
}

func (s *TokenService) Regenerate(id string) (*models.ApiToken, string, error) {
	fullToken, prefix, err := GenerateToken()
	if err != nil {
		return nil, "", err
	}

	newHash := HashToken(fullToken)

	token, err := s.repo.Regenerate(id, newHash, prefix)
	if err != nil {
		return nil, "", err
	}

	return token, fullToken, nil
}

func (s *TokenService) Validate(rawToken string) (*models.ApiToken, error) {
	hash := HashToken(rawToken)

	token, err := s.repo.FindByTokenHash(hash)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateLastUsed(token.ID); err != nil {
		logging.Warnf("failed to update token last_used_at for %s: %v", token.ID, err)
	}

	return token, nil
}

func generateUUID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[0:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:]), nil
}
