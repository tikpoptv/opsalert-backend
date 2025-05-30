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

func (s *Service) Update(ctx context.Context, id int, req *lineOAModel.UpdateRequest) error {
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

func (s *Service) List(ctx context.Context) ([]*lineOAModel.LineOA, error) {
	return s.repo.GetAll(ctx)
}
