package staff

type PermissionRequest struct {
	StaffID     int            `json:"staff_id" binding:"required"`
	Permissions []OAPermission `json:"permissions" binding:"required"`
}

type OAPermission struct {
	OAID            int    `json:"oa_id" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required,oneof=view manage"`
}
