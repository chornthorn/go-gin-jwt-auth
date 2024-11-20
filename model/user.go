package model

import "time"

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type User struct {
	ID        uint       `gorm:"primarykey"`
	Name      string     `json:"name" gorm:"not null"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Role      UserRole   `json:"role" gorm:"type:user_role;default:'user'"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
