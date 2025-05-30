package api_token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	apiTokenModel "opsalert/internal/model/api_token"
	apiTokenRepo "opsalert/internal/repository/api_token"
)

var (
	ErrUnauthorized = errors.New("unauthorized to create token for this OA")
)

type Service struct {
	repo apiTokenRepo.Repository
}

func NewService(repo apiTokenRepo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID int, oaID int, name string) (*apiTokenModel.APIToken, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	apiToken := &apiTokenModel.APIToken{
		UserID:    userID,
		OAID:      oaID,
		Token:     token,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, apiToken); err != nil {
		return nil, err
	}

	return apiToken, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*apiTokenModel.APIToken, error) {
	return s.repo.GetByToken(ctx, token)
}

func (s *Service) UpdateLastUsed(ctx context.Context, id int) error {
	return s.repo.UpdateLastUsed(ctx, id)
}

func (s *Service) CheckStaffOAPermission(ctx context.Context, staffID, oaID int) (bool, error) {
	return s.repo.CheckStaffOAPermission(ctx, staffID, oaID)
}
