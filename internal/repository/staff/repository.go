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

func (r *Repository) GetAll() ([]staffModel.Staff, error) {
	query := `
		SELECT id, username, password_hash, full_name, role, is_active, created_at
		FROM staff_accounts
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffs []staffModel.Staff
	for rows.Next() {
		var staff staffModel.Staff
		err := rows.Scan(
			&staff.ID,
			&staff.Username,
			&staff.PasswordHash,
			&staff.FullName,
			&staff.Role,
			&staff.IsActive,
			&staff.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		staffs = append(staffs, staff)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return staffs, nil
}

func (r *Repository) Update(id uint, staff *staffModel.Staff) error {
	query := `
		UPDATE staff_accounts
		SET full_name = $1, role = $2, is_active = $3
		WHERE id = $4`

	result, err := r.db.Exec(query, staff.FullName, staff.Role, staff.IsActive, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
