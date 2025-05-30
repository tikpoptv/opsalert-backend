package line_oa

import "time"

// LineOA represents a Line Official Account
type LineOA struct {
	ID                 uint      `json:"id" db:"id"`
	Name               string    `json:"name" db:"name" binding:"required,max=100"`
	ChannelID          string    `json:"channel_id" db:"channel_id" binding:"required,max=100"`
	ChannelSecret      string    `json:"channel_secret" db:"channel_secret" binding:"required"`
	ChannelAccessToken string    `json:"channel_access_token" db:"channel_access_token" binding:"required"`
	WebhookURL         string    `json:"webhook_url" db:"webhook_url"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type CreateRequest struct {
	Name               string `json:"name" binding:"required,max=100"`
	ChannelID          string `json:"channel_id" binding:"required,max=100"`
	ChannelSecret      string `json:"channel_secret" binding:"required"`
	ChannelAccessToken string `json:"channel_access_token" binding:"required"`
	WebhookURL         string `json:"webhook_url"`
}

type UpdateRequest struct {
	Name               string `json:"name" binding:"required,max=100"`
	ChannelID          string `json:"channel_id" binding:"required,max=100"`
	ChannelSecret      string `json:"channel_secret" binding:"required"`
	ChannelAccessToken string `json:"channel_access_token" binding:"required"`
	WebhookURL         string `json:"webhook_url"`
}
