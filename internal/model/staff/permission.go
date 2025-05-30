package staff

type PermissionRequest struct {
	StaffID     int            `json:"staff_id" binding:"required"`
	Permissions []OAPermission `json:"permissions" binding:"required"`
}

type OAPermission struct {
	OAID            int    `json:"oa_id" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required,oneof=view manage"`
}

type StaffPermissionResponse struct {
	OAID            int    `json:"oa_id"`
	OAName          string `json:"oa_name"`
	PermissionLevel string `json:"permission_level"`
}
