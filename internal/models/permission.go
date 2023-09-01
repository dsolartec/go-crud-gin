package models

type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"permission_name"`
	Description string `json:"description"`
	Deletable   bool   `json:"-"`
}

type UserPermission struct {
	UserID       int `json:"user_id"`
	PermissionID int `json:"permission_id"`
}
