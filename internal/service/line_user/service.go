package line_user

import (
	"context"
	"errors"
	"fmt"

	lineUserModel "opsalert/internal/model/line_user"
	lineUserRepo "opsalert/internal/repository/line_user"
	staffRepo "opsalert/internal/repository/staff"
)

var (
	ErrInsufficientPermissions = errors.New("insufficient permissions to view this OA")
)

type Service struct {
	repo      lineUserRepo.Repository
	staffRepo staffRepo.Repository
}

func NewService(repo lineUserRepo.Repository, staffRepo staffRepo.Repository) *Service {
	return &Service{
		repo:      repo,
		staffRepo: staffRepo,
	}
}

func (s *Service) GetByOaID(ctx context.Context, oaID int, staffID int, role string) ([]lineUserModel.LineUser, error) {
	// ถ้าไม่ใช่ admin ต้องตรวจสอบสิทธิ์
	if role != "admin" {
		hasPermission, err := s.staffRepo.CheckPermission(ctx, staffID, oaID)
		if err != nil {
			return nil, err
		}
		if !hasPermission {
			return nil, ErrInsufficientPermissions
		}
	}

	return s.repo.GetByOaID(ctx, oaID)
}

func (s *Service) GetByID(ctx context.Context, id uint, staffID int, role string) (*lineUserModel.LineUser, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// ถ้าเป็น admin สามารถดูข้อมูลได้เลย
	if role == "admin" {
		return user, nil
	}

	// ถ้าเป็น staff ต้องตรวจสอบสิทธิ์
	hasPermission, err := s.staffRepo.CheckPermission(ctx, staffID, int(user.OaID))
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions to view this OA")
	}

	return user, nil
}
