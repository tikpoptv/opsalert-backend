package line_user

import (
	"context"
	"database/sql"
	"fmt"
	lineUserModel "opsalert/internal/model/line_user"
)

type Repository interface {
	GetByOaID(ctx context.Context, oaID int) ([]lineUserModel.LineUser, error)
	GetByID(ctx context.Context, id uint) (*lineUserModel.LineUser, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByOaID(ctx context.Context, oaID int) ([]lineUserModel.LineUser, error) {
	query := `
		SELECT id, line_user_id, display_name, oa_id, created_at
		FROM line_users
		WHERE oa_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, oaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get line users: %w", err)
	}
	defer rows.Close()

	var users []lineUserModel.LineUser
	for rows.Next() {
		var user lineUserModel.LineUser
		if err := rows.Scan(&user.ID, &user.LineUserID, &user.DisplayName, &user.OaID, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan line user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating line users: %w", err)
	}

	return users, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*lineUserModel.LineUser, error) {
	query := `
		SELECT id, line_user_id, display_name, oa_id, created_at
		FROM line_users
		WHERE id = $1`

	var user lineUserModel.LineUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.LineUserID,
		&user.DisplayName,
		&user.OaID,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get line user: %w", err)
	}

	return &user, nil
}
