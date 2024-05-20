// Code generated by "stringer -type=Role -linecomment"; DO NOT EDIT.

package vobj

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RoleAll-0]
	_ = x[RoleSuperAdmin-1]
	_ = x[RoleAdmin-2]
	_ = x[RoleUser-3]
}

const _Role_name = "未知超级管理员管理员普通用户"

var _Role_index = [...]uint8{0, 6, 21, 30, 42}

func (i Role) String() string {
	if i < 0 || i >= Role(len(_Role_index)-1) {
		return "Role(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Role_name[_Role_index[i]:_Role_index[i+1]]
}

// IsAll 是否是：未知
func (i Role) IsAll() bool {
	return i == RoleAll
}

// IsSuperadmin 是否是：超级管理员
func (i Role) IsSuperadmin() bool {
	return i == RoleSuperAdmin
}

// IsAdmin 是否是：管理员
func (i Role) IsAdmin() bool {
	return i == RoleAdmin || i == RoleSuperAdmin
}

// IsUser 是否是：普通用户
func (i Role) IsUser() bool {
	return i == RoleUser
}

// GetValue 获取原始类型值
func (i Role) GetValue() int {
	return int(i)
}
