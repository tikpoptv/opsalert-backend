package line_oa

import (
	"context"
	"fmt"
	lineOAModel "opsalert/internal/model/line_oa"
	"opsalert/internal/repository/line_oa"
)

type Service struct {
	repo   *line_oa.Repository
	domain string
}

func NewService(repo *line_oa.Repository, domain string) *Service {
	return &Service{
		repo:   repo,
		domain: domain,
	}
}

func (s *Service) Create(ctx context.Context, req *lineOAModel.CreateRequest) error {
	webhookURL := fmt.Sprintf("https://%s/api/v1/webhook/line/%s", s.domain, req.ChannelID)

	oa := &lineOAModel.LineOA{
		Name:               req.Name,
		ChannelID:          req.ChannelID,
		ChannelSecret:      req.ChannelSecret,
		ChannelAccessToken: req.ChannelAccessToken,
		WebhookURL:         webhookURL,
	}

	return s.repo.Create(ctx, oa)
}

func (s *Service) Update(ctx context.Context, id int, staffID int, role string, req *lineOAModel.UpdateRequest) error {
	// ถ้าเป็น admin ให้แก้ไขได้เลย
	if role == "admin" {
		oa, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if oa == nil {
			return fmt.Errorf("line official account not found")
		}

		webhookURL := fmt.Sprintf("https://%s/api/v1/webhook/line/%s", s.domain, req.ChannelID)

		oa.Name = req.Name
		oa.ChannelID = req.ChannelID
		oa.ChannelSecret = req.ChannelSecret
		oa.ChannelAccessToken = req.ChannelAccessToken
		oa.WebhookURL = webhookURL

		return s.repo.Update(ctx, oa)
	}

	// ถ้าเป็น staff ต้องเช็คสิทธิ์ manage
	hasPermission, err := s.repo.CheckManagePermission(ctx, staffID, id)
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("insufficient permissions")
	}

	// ถ้ามีสิทธิ์ manage ให้แก้ไขได้
	oa, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if oa == nil {
		return fmt.Errorf("line official account not found")
	}

	webhookURL := fmt.Sprintf("https://%s/api/v1/webhook/line/%s", s.domain, req.ChannelID)

	oa.Name = req.Name
	oa.ChannelID = req.ChannelID
	oa.ChannelSecret = req.ChannelSecret
	oa.ChannelAccessToken = req.ChannelAccessToken
	oa.WebhookURL = webhookURL

	return s.repo.Update(ctx, oa)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	oa, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if oa == nil {
		return fmt.Errorf("line official account not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, staffID int, role string) ([]*lineOAModel.LineOA, error) {
	// ถ้าเป็น admin ให้แสดงทั้งหมด
	if role == "admin" {
		return s.repo.GetAll(ctx)
	}
	// ถ้าเป็น staff ให้แสดงเฉพาะที่ได้รับสิทธิ์
	return s.repo.GetByStaffID(ctx, staffID)
}
