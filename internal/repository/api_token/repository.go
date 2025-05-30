package api_token

import (
	"context"
	"database/sql"
	"time"

	"opsalert/internal/model/api_token"
)

type Repository interface {
	Create(ctx context.Context, token *api_token.APIToken) error
	GetByToken(ctx context.Context, token string) (*api_token.APIToken, error)
	UpdateLastUsed(ctx context.Context, id int) error
	CheckStaffOAPermission(ctx context.Context, staffID, oaID int) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, token *api_token.APIToken) error {
	query := `
		INSERT INTO api_tokens (user_id, oa_id, token, name, created_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.OAID,
		token.Token,
		token.Name,
		time.Now(),
		time.Now(),
	).Scan(&token.ID)
}

func (r *repository) GetByToken(ctx context.Context, token string) (*api_token.APIToken, error) {
	query := `
		SELECT id, user_id, oa_id, token, name, created_at, last_used_at
		FROM api_tokens
		WHERE token = $1`

	var t api_token.APIToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.ID,
		&t.UserID,
		&t.OAID,
		&t.Token,
		&t.Name,
		&t.CreatedAt,
		&t.LastUsedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) UpdateLastUsed(ctx context.Context, id int) error {
	query := `
		UPDATE api_tokens
		SET last_used_at = $1
		WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *repository) CheckStaffOAPermission(ctx context.Context, staffID, oaID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM staff_oa_permissions
			WHERE staff_id = $1 AND oa_id = $2
		)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, staffID, oaID).Scan(&exists)
	return exists, err
}
