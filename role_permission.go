package authority

type RolePermissions []*RolePermission

var NilRolePermission = (*RolePermission)(nil)

type RolePermission struct {
	RoleID       uint64 `json:"role_id,omitempty" gorm:"primaryKey"`
	PermissionID uint64 `json:"permission_id,omitempty" gorm:"primaryKey"`
}

func (RolePermission) TableName() string {
	return auth.prefix + "role_permissions"
}
