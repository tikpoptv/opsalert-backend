package staff

import (
	"database/sql"
	staffModel "opsalert/internal/model/staff"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(staff *staffModel.Staff) error {
	query := `
		INSERT INTO staff_accounts (username, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		staff.Username,
		staff.PasswordHash,
		staff.FullName,
		staff.Role,
		staff.IsActive,
	).Scan(&staff.ID, &staff.CreatedAt)
}

func (r *Repository) GetByUsername(username string) (*staffModel.Staff, error) {
	query := `
		SELECT id, username, password_hash, full_name, role, is_active, created_at
		FROM staff_accounts
		WHERE username = $1`

	staff := &staffModel.Staff{}
	err := r.db.QueryRow(query, username).Scan(
		&staff.ID,
		&staff.Username,
		&staff.PasswordHash,
		&staff.FullName,
		&staff.Role,
		&staff.IsActive,
		&staff.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *Repository) GetByID(id uint) (*staffModel.Staff, error) {
	query := `
		SELECT id, username, password_hash, full_name, role, is_active, created_at
		FROM staff_accounts
		WHERE id = $1`

	staff := &staffModel.Staff{}
	err := r.db.QueryRow(query, id).Scan(
		&staff.ID,
		&staff.Username,
		&staff.PasswordHash,
		&staff.FullName,
		&staff.Role,
		&staff.IsActive,
		&staff.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return staff, nil
}
