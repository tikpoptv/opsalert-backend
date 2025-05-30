package line_oa

import "time"

type LineOA struct {
	ID                 uint      `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	ChannelID          string    `json:"channel_id" db:"channel_id"`
	ChannelSecret      string    `json:"channel_secret" db:"channel_secret"`
	ChannelAccessToken string    `json:"channel_access_token" db:"channel_access_token"`
	WebhookURL         string    `json:"webhook_url" db:"webhook_url"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
