package authority

type UserRoles []*UserRole

var NilUserRole = (*UserRole)(nil)

type UserRole struct {
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	RoleID uint64 `json:"role_id" gorm:"primaryKey"`
}

func (UserRole) TableName() string {
	return auth.prefix + "user_roles"
}
