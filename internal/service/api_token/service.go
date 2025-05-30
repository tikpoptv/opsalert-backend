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

func (s *Service) Create(ctx context.Context, userID int, name string) (*apiTokenModel.APIToken, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	apiToken := &apiTokenModel.APIToken{
		UserID:    userID,
		Token:     token,
		Name:      name,
		IsActive:  true,
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

func (s *Service) ResetToken(ctx context.Context, tokenID int, userID int, role string) (*apiTokenModel.APIToken, error) {
	// ตรวจสอบว่า token นี้เป็นของ user นี้หรือไม่
	token, err := s.repo.GetByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	// ถ้าไม่ใช่ admin ต้องเป็นเจ้าของ token
	if role != "admin" && token.UserID != userID {
		return nil, ErrUnauthorized
	}

	// สร้าง token ใหม่
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	newToken := hex.EncodeToString(tokenBytes)

	// อัพเดท token ในฐานข้อมูล
	token.Token = newToken
	token.LastUsedAt = time.Now()
	if err := s.repo.Update(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) UpdateStatus(ctx context.Context, tokenID int, userID int, role string, isActive bool) (*apiTokenModel.APIToken, error) {
	// ตรวจสอบว่า token นี้เป็นของ user นี้หรือไม่
	token, err := s.repo.GetByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	// ถ้าไม่ใช่ admin ต้องเป็นเจ้าของ token
	if role != "admin" && token.UserID != userID {
		return nil, ErrUnauthorized
	}

	// อัพเดทสถานะ token
	token.IsActive = isActive
	if err := s.repo.Update(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}
