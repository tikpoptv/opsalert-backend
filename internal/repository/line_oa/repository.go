package line_oa

import (
	"context"
	"database/sql"
	"fmt"
	lineOAModel "opsalert/internal/model/line_oa"

	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, oa *lineOAModel.LineOA) error
	GetByID(ctx context.Context, id int) (*lineOAModel.LineOA, error)
	Update(ctx context.Context, oa *lineOAModel.LineOA) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*lineOAModel.LineOA, error)
	GetByStaffID(ctx context.Context, staffID int) ([]*lineOAModel.LineOA, error)
	CheckManagePermission(ctx context.Context, staffID int, oaID int) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, oa *lineOAModel.LineOA) error {
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

func (r *repository) GetByID(ctx context.Context, id int) (*lineOAModel.LineOA, error) {
	query := `
		SELECT id, name, channel_id, channel_secret, channel_access_token, webhook_url, created_at
		FROM line_official_accounts
		WHERE id = $1`

	var oa lineOAModel.LineOA
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
		return nil, fmt.Errorf("failed to get OA: %w", err)
	}

	return &oa, nil
}

func (r *repository) Update(ctx context.Context, oa *lineOAModel.LineOA) error {
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

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM line_official_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]*lineOAModel.LineOA, error) {
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

func (r *repository) GetByStaffID(ctx context.Context, staffID int) ([]*lineOAModel.LineOA, error) {
	// 1. เช็คสิทธิ์ก่อน
	var oaIDs []int
	rows, err := r.db.QueryContext(ctx, "SELECT oa_id FROM staff_oa_permissions WHERE staff_id = $1", staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oaID int
		if err := rows.Scan(&oaID); err != nil {
			return nil, fmt.Errorf("failed to scan oa_id: %w", err)
		}
		oaIDs = append(oaIDs, oaID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	// ถ้าไม่มีสิทธิ์เลย
	if len(oaIDs) == 0 {
		return make([]*lineOAModel.LineOA, 0), nil
	}

	// 2. ดึงข้อมูล OA จาก oa_id ที่มีสิทธิ์
	query := `
		SELECT id, name, channel_id, channel_secret, channel_access_token, webhook_url, created_at
		FROM line_official_accounts
		WHERE id = ANY($1)
		ORDER BY created_at DESC`

	rows, err = r.db.QueryContext(ctx, query, pq.Array(oaIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query OAs: %w", err)
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
			return nil, fmt.Errorf("failed to scan OA: %w", err)
		}
		oas = append(oas, oa)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return oas, nil
}

func (r *repository) CheckManagePermission(ctx context.Context, staffID int, oaID int) (bool, error) {
	var hasPermission bool
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM staff_oa_permissions 
			WHERE staff_id = $1 
			AND oa_id = $2 
			AND permission_level = 'manage'
		)`

	err := r.db.QueryRowContext(ctx, query, staffID, oaID).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check manage permission: %w", err)
	}

	return hasPermission, nil
}
