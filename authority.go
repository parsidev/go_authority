package authority

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Authority struct {
	prefix string
	db     *gorm.DB
	opts   *Options
}

var auth *Authority

func New(opts ...Option) (a *Authority, err error) {
	options := newOptions(opts...)

	a = &Authority{
		opts:   options,
		prefix: options.prefix,
		db:     options.db,
	}

	auth = a

	if options.migrate.Valid && (options.migrate.Bool == true) {
		if err = a.db.AutoMigrate(&Role{}, &Permission{}, &RolePermission{}, &UserRole{}); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func Resolve() *Authority {
	return auth
}

func (a *Authority) CreateRole(r *Role) (err error) {
	var (
		role  = &Role{}
		rName = r.Name
	)

	if err = a.db.First(&role, "name = ?", rName).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil && role.IsValid() {
		return errors.New(fmt.Sprintf("role '%v' already exists", rName))

	}

	if err = a.db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (a *Authority) CreatePermission(p *Permission) (err error) {
	var (
		permission = &Permission{}
		pName      = p.Name
	)

	if err = a.db.First(&permission, "name = ?", pName).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil && permission.IsValid() {
		return errors.New(fmt.Sprintf("permission '%v' already exists", pName))
	}

	if err = a.db.Create(&p).Error; err != nil {
		return err
	}

	return nil
}

func (a *Authority) AssignPermissionsToRole(req *RolePermissionRequest) (err error) {
	var (
		rolePermissions = make(RolePermissions, 0)
		rolePermission  *RolePermission
		permission      = new(Permission)
		ok              bool
		role            = &Role{}
	)

	if err = a.db.First(&role, "id = ?", req.RoleID).Error; err != nil {
		return err
	}

	for _, permID := range req.PermissionIDs {
		if err = a.db.First(&permission, "id = ?", permID).Error; err != nil {
			continue
		}

		rolePermission = &RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
		}

		if ok, err = a.CheckRolePermission(rolePermission); ok || err != nil {
			continue
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	if len(rolePermissions) == 0 {
		return errors.New(fmt.Sprintf("all permissions already assigned to role '%s'", role.Name))
	}

	if err = a.db.Create(&rolePermissions).Error; err != nil {
		return err
	}

	return nil

}

func (a *Authority) AssignRoleToUser(req *UserRole) (err error) {
	var (
		role = &Role{}
		ok   bool
	)

	if err = a.db.First(&role, "id = ?", req.RoleID).Error; err != nil {
		return err
	}

	if ok, err = a.CheckUserRole(req); err != nil {
		return err
	}

	if ok {
		return errors.New(fmt.Sprintf("role '%s' already assigned to user '%v'", role.Name, req.UserID))
	}

	if err = a.db.Create(&req).Error; err != nil {
		return err
	}

	return nil

}

func (a *Authority) CheckUserRole(req *UserRole) (ok bool, err error) {
	if err = a.db.
		Model(&UserRole{}).
		Select("count(*) > 0").
		Where("role_id = ?", req.RoleID).
		Where("user_id = ?", req.UserID).Find(&ok).Error; err != nil {
		return false, err
	}

	return ok, nil
}

func (a *Authority) CheckUserPermission(userID uint64, permName string) (ok bool, err error) {
	var (
		userRoles = make([]*UserRole, 0)
		roleIDs   = make([]uint64, 0)
		perm      *Permission
	)

	if err = a.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrUserDontHaveRole
		}
		return false, err
	}

	for _, r := range userRoles {
		roleIDs = append(roleIDs, r.RoleID)
	}

	if err = a.db.First(&perm, "name = ?", permName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrPermissionNotFound
		}
		return false, err
	}

	if err = a.db.
		Model(&RolePermission{}).
		Select("count(*) > 0").
		Where("role_id in (?)", roleIDs).
		Where("permission_id = ?", perm.ID).Find(&ok).Error; err != nil {
		return false, err
	}

	return ok, nil
}

func (a *Authority) CheckRolePermission(rp *RolePermission) (ok bool, err error) {
	if err = a.db.
		Model(&RolePermission{}).
		Select("count(*) > 0").
		Where("role_id = ?", rp.RoleID).
		Where("permission_id = ?", rp.PermissionID).Find(&ok).Error; err != nil {
		return false, err
	}

	return ok, nil
}

func (a *Authority) RevokeUserRole(userID, roleID uint64) (err error) {
	var (
		ok   bool
		role = &Role{}
	)

	if ok, err = a.CheckUserRole(&UserRole{RoleID: roleID, UserID: userID}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if !ok {
		if err = a.db.First(&role, "id = ?", roleID).Error; err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("role '%s' not assigned to user '%v'", role.Name, userID))
	}

	if err = a.db.
		Where("role_id = ?", roleID).
		Where("user_id = ?", userID).
		Delete(&UserRole{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (a *Authority) RevokeRolePermission(roleID, permID uint64) (err error) {
	var (
		ok         bool
		role       = &Role{}
		permission = &Permission{}
	)

	if ok, err = a.CheckRolePermission(&RolePermission{RoleID: roleID, PermissionID: permID}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if !ok {

		if err = a.db.First(&role, "id = ?", roleID).Error; err != nil {
			return err
		}

		if err = a.db.First(&permission, "permission_id = ?", permID).Error; err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("permission '%s' not assigned to role '%s'", permission.Name, role.Name))
	}

	if err = a.db.
		Where("role_id = ?", roleID).
		Where("permission_id = ?", permID).
		Delete(&RolePermission{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (a *Authority) GetUserRoles(userID uint64) (roles Roles, err error) {
	var (
		userRoles = make(UserRoles, 0)
		roleIDs   = make([]uint64, 0)
	)

	roles = make(Roles, 0)

	if err = a.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	for _, uR := range userRoles {
		roleIDs = append(roleIDs, uR.RoleID)
	}

	if a.db.Where("id IN (?)", roleIDs).Find(&roles).Error != nil {
		return nil, err
	}

	return roles, nil
}

func (a *Authority) GetRolePermissions(roleID uint64) (permissions Permissions, err error) {
	var (
		rolePerms = make(RolePermissions, 0)
		permIDs   = make([]uint64, 0)
	)

	permissions = make(Permissions, 0)

	if err = a.db.Where("role_id = ?", roleID).Find(&rolePerms).Error; err != nil {
		return nil, err
	}

	for _, p := range rolePerms {
		permIDs = append(permIDs, p.PermissionID)
	}

	if a.db.Where("id IN (?)", permIDs).Find(&permissions).Error != nil {
		return nil, err
	}

	return permissions, nil
}

func (a *Authority) FetchAllRoles() (roles Roles, err error) {
	roles = make(Roles, 0)

	if err = a.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (a *Authority) FetchAllPermissions() (permissions Permissions, err error) {
	permissions = make(Permissions, 0)

	if err = a.db.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (a *Authority) DeleteRole(roleID uint64) (err error) {
	var (
		role = &Role{}
		c    int64
	)

	if err = a.db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}

		return err
	}

	if err = a.db.Model(&UserRole{}).Where("role_id = ?", roleID).Count(&c).Error; err != nil {
		return err
	}

	if c > 0 {
		return ErrRoleInUse
	}

	tx := a.db.Begin()

	if err = tx.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Delete(&Role{}, "id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (a *Authority) DeletePermission(permID uint64) (err error) {
	var (
		permission = &Permission{}
		c          int64
	)

	if err = a.db.First(&permission, "id = ?", permID).Error; err != nil {
		return err
	}

	if err = a.db.Model(&RolePermission{}).Where("permission_id = ?", permID).Count(&c).Error; err != nil {
		return err
	}

	if c > 0 {
		return ErrPermissionInUse
	}

	if err = a.db.Delete(&Permission{}, "id = ?", permID).Error; err != nil {
		return err
	}

	return nil
}

func (a *Authority) GetPermission(permName string) (permission *Permission, err error) {
	permission = new(Permission)

	if err = a.db.Where("name = ?", permName).First(&permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (a *Authority) GetRole(roleName string) (role *Role, err error) {
	role = new(Role)

	if err = a.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}

	return role, nil
}
