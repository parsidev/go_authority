package authority

import "errors"

var (
	ErrPermissionInUse    = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrRoleInUse          = errors.New("cannot delete assigned role")
	ErrRoleNotFound       = errors.New("role not found")
	ErrUserDontHaveRole   = errors.New("user don't have role")
)
