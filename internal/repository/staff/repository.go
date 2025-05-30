package staff

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *Repository) SetPermissions(ctx context.Context, staffID int, permissions []staffModel.OAPermission) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// ตรวจสอบว่า staff มีอยู่จริง
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM staff_accounts WHERE id = $1)", staffID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check staff existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("staff not found")
	}

	// ตรวจสอบว่า OA ทั้งหมดมีอยู่จริง
	for _, p := range permissions {
		err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM line_official_accounts WHERE id = $1)", p.OAID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check OA existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("OA with ID %d not found", p.OAID)
		}
	}

	// ลบสิทธิ์เดิมทั้งหมด
	_, err = tx.ExecContext(ctx, "DELETE FROM staff_oa_permissions WHERE staff_id = $1", staffID)
	if err != nil {
		return fmt.Errorf("failed to delete existing permissions: %w", err)
	}

	// เพิ่มสิทธิ์ใหม่
	for _, p := range permissions {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO staff_oa_permissions (staff_id, oa_id, permission_level) VALUES ($1, $2, $3)",
			staffID, p.OAID, p.PermissionLevel,
		)
		if err != nil {
			return fmt.Errorf("failed to insert permission for OA %d: %w", p.OAID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetStaffPermissions(ctx context.Context, staffID int) ([]staffModel.StaffPermissionResponse, error) {
	query := `
		SELECT p.oa_id, oa.name, p.permission_level
		FROM staff_oa_permissions p
		JOIN line_official_accounts oa ON p.oa_id = oa.id
		WHERE p.staff_id = $1
		ORDER BY oa.name`

	rows, err := r.db.QueryContext(ctx, query, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff permissions: %w", err)
	}
	defer rows.Close()

	var permissions []staffModel.StaffPermissionResponse
	for rows.Next() {
		var p staffModel.StaffPermissionResponse
		if err := rows.Scan(&p.OAID, &p.OAName, &p.PermissionLevel); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

func (r *Repository) DeleteStaffPermissions(ctx context.Context, staffID int, oaID int) error {
	// ตรวจสอบว่า staff มีอยู่จริง
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM staff_accounts WHERE id = $1)", staffID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check staff existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("staff not found")
	}

	// ตรวจสอบว่า OA มีอยู่จริง
	err = r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM line_official_accounts WHERE id = $1)", oaID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check OA existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("OA not found")
	}

	// ลบสิทธิ์เฉพาะ OA ที่ระบุ
	result, err := r.db.ExecContext(ctx, "DELETE FROM staff_oa_permissions WHERE staff_id = $1 AND oa_id = $2", staffID, oaID)
	if err != nil {
		return fmt.Errorf("failed to delete staff permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("staff does not have permission for this OA")
	}

	return nil
}
