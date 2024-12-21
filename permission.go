package authority

type Permissions []*Permission

var NilPermission = (*Permission)(nil)

type Permission struct {
	Base
	Name        string `json:"name" query:"name" param:"name" gorm:"size:50;not null;unique"`
	DisplayName string `json:"display_name" gorm:"size:50;null"`
}

func (Permission) TableName() string {
	return auth.prefix + "permissions"
}

type PermissionRequest struct {
	RequestData
	Permission
}
