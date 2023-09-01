package services

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/models"
	"go-crud-gin/internal/platform/logger"
	"strings"
)

type PermissionsService interface {
	Create(name, description string) (int, error)
	GetPermissionByID(id int) *models.Permission
	GetPermissionByName(name string) *models.Permission
	GetPermissions() []models.Permission
	DeletePermission(name string) error
	GetPermissionsForUser(userID int) []models.UserPermission
	UserHasPermission(userID, permissionID int) bool
	GrantPermissionToUser(userID int, permissionName string) error
	RevokePermissionToUser(userID int, permissionName string) error
}

type permissionsService struct {
	BaseService

	permissions     []models.Permission
	userPermissions []models.UserPermission
}

func (service *permissionsService) Create(name, description string) (int, error) {
	lastID := 0
	for _, permission := range service.permissions {
		if strings.EqualFold(permission.Name, name) {
			return 0, apperror.NewErrPermissionAlreadyExists()
		}

		if permission.ID > lastID {
			lastID = permission.ID
		}
	}

	lastID++

	service.permissions = append(service.permissions, models.Permission{
		ID:          lastID,
		Name:        name,
		Description: description,
		Deletable:   true,
	})

	service.logger.Infof("[PermissionsService] New permission created %s!", name)

	return lastID, nil
}

func (service *permissionsService) GetPermissionByID(id int) *models.Permission {
	for _, permission := range service.permissions {
		if permission.ID == id {
			return &permission
		}
	}

	return nil
}

func (service *permissionsService) GetPermissionByName(name string) *models.Permission {
	for _, permission := range service.permissions {
		if strings.EqualFold(permission.Name, name) {
			return &permission
		}
	}

	return nil
}

func (service *permissionsService) GetPermissions() []models.Permission {
	return service.permissions
}

func (service *permissionsService) DeletePermission(name string) error {
	var permissionIDDeleted *int

	newPermissions := []models.Permission{}
	for _, permission := range service.permissions {
		if strings.EqualFold(permission.Name, name) {
			if !permission.Deletable {
				return apperror.NewErrPermissionNotDeletable()
			}

			permissionIDDeleted = &permission.ID

			continue
		}

		newPermissions = append(newPermissions, permission)
	}

	if permissionIDDeleted == nil {
		return apperror.NewErrPermissionNotFound()
	}

	service.permissions = newPermissions

	newUserPermissions := []models.UserPermission{}
	for _, userPermission := range service.userPermissions {
		if userPermission.PermissionID == *permissionIDDeleted {
			continue
		}

		newUserPermissions = append(newUserPermissions, userPermission)
	}

	service.userPermissions = newUserPermissions

	return nil
}

func (service *permissionsService) GetPermissionsForUser(userID int) []models.UserPermission {
	permissions := []models.UserPermission{}
	for _, userPermission := range service.userPermissions {
		if userPermission.UserID == userID {
			permissions = append(permissions, userPermission)
		}
	}

	return permissions
}

func (service *permissionsService) UserHasPermission(userID, permissionID int) bool {
	for _, userPermission := range service.userPermissions {
		if userPermission.UserID == userID && userPermission.PermissionID == permissionID {
			return true
		}
	}

	return false
}

func (service *permissionsService) GrantPermissionToUser(userID int, permissionName string) error {
	permission := service.GetPermissionByName(permissionName)
	if permission == nil {
		return apperror.NewErrPermissionNotFound()
	}

	hasPermission := service.UserHasPermission(userID, permission.ID)
	if hasPermission {
		return apperror.NewErrUserAlreadyHasPermission()
	}

	service.userPermissions = append(service.userPermissions, models.UserPermission{
		UserID:       userID,
		PermissionID: permission.ID,
	})

	service.logger.Infof("[PermissionsService] Permission '%s' granted to '%s' user!", permissionName, userID)

	return nil
}

func (service *permissionsService) RevokePermissionToUser(userID int, permissionName string) error {
	permission := service.GetPermissionByName(permissionName)
	if permission == nil {
		return apperror.NewErrPermissionNotFound()
	}

	newPermissions := []models.UserPermission{}
	for i := 0; i < len(service.userPermissions); i++ {
		value := service.userPermissions[i]
		if value.PermissionID == permission.ID && value.UserID == userID {
			continue
		}

		newPermissions = append(newPermissions, value)
	}

	if len(service.userPermissions) == len(newPermissions) {
		return apperror.NewErrUserPermissionNotFound()
	}

	service.userPermissions = newPermissions

	return nil
}

func NewPermissionsService(
	logger logger.Logger,
) PermissionsService {
	return &permissionsService{
		BaseService: BaseService{
			logger: logger,
		},

		permissions: []models.Permission{
			{
				ID:          1,
				Name:        "users_full",
				Description: "Full access to users endpoints",
				Deletable:   false,
			},
			{
				ID:          2,
				Name:        "users_read",
				Description: "Only access to users GET endpoints",
				Deletable:   false,
			},
			{
				ID:          3,
				Name:        "users_write",
				Description: "Only access to users POST, PUT and DELETE endpoints",
				Deletable:   false,
			},
			{
				ID:          4,
				Name:        "permissions_full",
				Description: "Full access to permissions endpoints",
				Deletable:   false,
			},
			{
				ID:          5,
				Name:        "permissions_read",
				Description: "Only access to permissions GET endpoints",
				Deletable:   false,
			},
			{
				ID:          6,
				Name:        "permissions_write",
				Description: "Only access to permissions POST, PUT and DELETE endpoints",
				Deletable:   false,
			},
			{
				ID:          7,
				Name:        "grant_permission",
				Description: "Grant a permission to an user",
				Deletable:   false,
			},
			{
				ID:          8,
				Name:        "revoke_permission",
				Description: "Revoke a permission to an user",
				Deletable:   false,
			},
		},
		userPermissions: []models.UserPermission{
			{
				UserID:       1,
				PermissionID: 1,
			},
			{
				UserID:       1,
				PermissionID: 4,
			},
			{
				UserID:       1,
				PermissionID: 7,
			},
			{
				UserID:       1,
				PermissionID: 8,
			},
			{
				UserID:       2,
				PermissionID: 2,
			},
			{
				UserID:       2,
				PermissionID: 5,
			},
		},
	}
}
