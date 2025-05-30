package line_oa

import (
	"database/sql"
	lineOAModel "opsalert/internal/model/line_oa"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(oa *lineOAModel.LineOA) error {
	query := `
		INSERT INTO line_official_accounts (name, channel_id, channel_secret, channel_access_token, webhook_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		oa.Name,
		oa.ChannelID,
		oa.ChannelSecret,
		oa.ChannelAccessToken,
		oa.WebhookURL,
	).Scan(&oa.ID, &oa.CreatedAt)
}
