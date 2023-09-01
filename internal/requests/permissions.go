package requests

type CreatePermission struct {
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
}
