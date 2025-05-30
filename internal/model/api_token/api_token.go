package api_token

import "time"

type APIToken struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	OAID       int       `json:"oa_id"`
	Token      string    `json:"token"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
	OAID int    `json:"oa_id" binding:"required"`
}

type CreateResponse struct {
	Token string `json:"token"`
}
