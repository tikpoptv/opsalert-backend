package line_user

import "time"

type LineUser struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	OaID      int       `json:"oa_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListResponse struct {
	Data []LineUser `json:"data"`
}
