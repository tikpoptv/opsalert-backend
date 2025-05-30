package staff

type PermissionRequest struct {
	StaffID     int            `json:"staff_id" binding:"required"`
	Permissions []OAPermission `json:"permissions" binding:"required"`
}

type OAPermission struct {
	ID              uint   `json:"id" db:"id"`
	StaffID         uint   `json:"staff_id" db:"staff_id" binding:"required"`
	OAID            uint   `json:"oa_id" db:"oa_id" binding:"required"`
	PermissionLevel string `json:"permission_level" db:"permission_level" binding:"required,oneof=view manage"`
}

type StaffPermissionResponse struct {
	OAID            int    `json:"oa_id"`
	OAName          string `json:"oa_name"`
	PermissionLevel string `json:"permission_level"`
}
