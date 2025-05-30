package line_oa

import (
	"fmt"
	"opsalert/config"
	lineOAModel "opsalert/internal/model/line_oa"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(oa *lineOAModel.LineOA) error
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOA(req *lineOAModel.CreateLineOARequest) error {
	// สร้าง webhook URL จาก domain ของเรา
	webhookURL := fmt.Sprintf("https://%s/api/v1/webhook/line/%s", config.Get().Domain, req.ChannelID)

	oa := &lineOAModel.LineOA{
		Name:               req.Name,
		ChannelID:          req.ChannelID,
		ChannelSecret:      req.ChannelSecret,
		ChannelAccessToken: req.ChannelAccessToken,
		WebhookURL:         webhookURL,
	}

	return s.repo.Create(oa)
}
