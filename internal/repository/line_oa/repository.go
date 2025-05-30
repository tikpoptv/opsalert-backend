package line_oa

import (
	"context"
	"database/sql"
	lineOAModel "opsalert/internal/model/line_oa"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, oa *lineOAModel.LineOA) error {
	query := `
		INSERT INTO line_official_accounts (name, channel_id, channel_secret, channel_access_token, webhook_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		oa.Name,
		oa.ChannelID,
		oa.ChannelSecret,
		oa.ChannelAccessToken,
		oa.WebhookURL,
	).Scan(&oa.ID, &oa.CreatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id int) (*lineOAModel.LineOA, error) {
	query := `
		SELECT id, name, channel_id, channel_secret, channel_access_token, webhook_url, created_at
		FROM line_official_accounts
		WHERE id = $1`

	oa := &lineOAModel.LineOA{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&oa.ID,
		&oa.Name,
		&oa.ChannelID,
		&oa.ChannelSecret,
		&oa.ChannelAccessToken,
		&oa.WebhookURL,
		&oa.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return oa, nil
}

func (r *Repository) Update(ctx context.Context, oa *lineOAModel.LineOA) error {
	query := `
		UPDATE line_official_accounts
		SET name = $1,
			channel_id = $2,
			channel_secret = $3,
			channel_access_token = $4,
			webhook_url = $5
		WHERE id = $6`

	_, err := r.db.ExecContext(
		ctx,
		query,
		oa.Name,
		oa.ChannelID,
		oa.ChannelSecret,
		oa.ChannelAccessToken,
		oa.WebhookURL,
		oa.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM line_official_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) GetAll(ctx context.Context) ([]*lineOAModel.LineOA, error) {
	query := `
		SELECT id, name, channel_id, channel_secret, channel_access_token, webhook_url, created_at
		FROM line_official_accounts
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oas []*lineOAModel.LineOA
	for rows.Next() {
		oa := &lineOAModel.LineOA{}
		err := rows.Scan(
			&oa.ID,
			&oa.Name,
			&oa.ChannelID,
			&oa.ChannelSecret,
			&oa.ChannelAccessToken,
			&oa.WebhookURL,
			&oa.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		oas = append(oas, oa)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return oas, nil
}
