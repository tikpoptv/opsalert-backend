package line_user

import "time"

type LineUser struct {
	ID          uint      `json:"id" db:"id"`
	LineUserID  string    `json:"line_user_id" db:"line_user_id" binding:"required,max=64"`
	DisplayName string    `json:"display_name" db:"display_name"`
	OaID        uint      `json:"oa_id" db:"oa_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ListResponse struct {
	Data []LineUser `json:"data"`
}
