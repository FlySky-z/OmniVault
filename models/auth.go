package models

// Role 代表系统中的角色
type Role struct {
	Name        string       `gorm:"unique;not null"`             // 角色名称
	Description string       `gorm:"type:varchar(255);"`          // 角色描述
	Permissions []Permission `gorm:"many2many:role_permissions;"` // 角色拥有的权限（多对多关系）
	UserRoles   []UserRole   `gorm:"foreignKey:RoleID"`           // 用户角色关联
	ID          uint         `gorm:"primaryKey"`
}

// Role 表名设置
func (Role) TableName() string {
	return "roles"
}

// Permission 代表系统中的权限
type Permission struct {
	Name        string `gorm:"unique;not null"`             // 权限名称
	Description string `gorm:"type:varchar(255);"`          // 权限描述
	Roles       []Role `gorm:"many2many:role_permissions;"` // 权限属于的角色（多对多关系）
	ID          uint   `gorm:"primaryKey"`
}

// Permission 表名设置
func (Permission) TableName() string {
	return "permissions"
}

// User 代表系统中的用户
type User struct {
	Username  string     `gorm:"unique;not null"`    // 用户名（唯一）
	Password  string     `gorm:"not null"`           // 加密后的密码
	Email     string     `gorm:"type:varchar(255);"` // 用户邮箱（可选）
	UserRoles []UserRole `gorm:"foreignKey:UserID"`  // 用户角色关联
	ID        uint       `gorm:"primaryKey"`
}

// User 表名设置
func (User) TableName() string {
	return "users"
}

// UserRole 代表用户与角色的关联
type UserRole struct {
	User   User `gorm:"foreignKey:UserID"`
	Role   Role `gorm:"foreignKey:RoleID"`
	UserID uint `gorm:"primaryKey"` // 用户ID
	RoleID uint `gorm:"primaryKey"` // 角色ID
}

// UserRole 表名设置
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 代表角色与权限的关联
type RolePermission struct {
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	RoleID       uint       `gorm:"primaryKey"` // 角色ID
	PermissionID uint       `gorm:"primaryKey"` // 权限ID
}

// RolePermission 表名设置
func (RolePermission) TableName() string {
	return "role_permissions"
}
