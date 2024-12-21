package authority

type Roles []*Role

var NilRole = (*Role)(nil)

type Role struct {
	Base
	Name        string `json:"name" query:"name" param:"name" gorm:"size:50;not null;unique"`
	DisplayName string `json:"display_name" gorm:"size:50;null"`
}

func (Role) TableName() string {
	return auth.prefix + "roles"
}

type RoleRequest struct {
	RequestData
	Role
}
