package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"opsalert/internal/repository/line_oa"
)

type Service struct {
	oaRepo     line_oa.Repository
	httpClient *http.Client
}

func NewService(oaRepo line_oa.Repository) *Service {
	return &Service{
		oaRepo: oaRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ForwardWebhook ส่ง webhook ไปยังระบบภายนอก
func (s *Service) ForwardWebhook(ctx context.Context, oaID int, data map[string]interface{}) error {
	// ดึงข้อมูล OA
	oa, err := s.oaRepo.GetByID(ctx, oaID)
	if err != nil {
		return fmt.Errorf("failed to get OA: %w", err)
	}

	if oa.WebhookURL == "" {
		return fmt.Errorf("no webhook URL configured for OA %d", oaID)
	}

	// แปลงข้อมูลเป็น JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook data: %w", err)
	}

	// ส่ง POST request
	req, err := http.NewRequestWithContext(ctx, "POST", oa.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook request failed with status %d", resp.StatusCode)
	}

	return nil
}
