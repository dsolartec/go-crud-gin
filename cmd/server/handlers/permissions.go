package handlers

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/models"
	"go-crud-gin/internal/platform/logger"
	"go-crud-gin/internal/requests"
	"go-crud-gin/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PermissionsHandler struct {
	BaseHandler

	permissionsService services.PermissionsService
	usersService       services.UsersService
}

func (handler *PermissionsHandler) CreatePermission(c *gin.Context) error {
	var body *requests.CreatePermission
	if err := c.ShouldBind(&body); err != nil {
		return err
	}

	validationErrors := map[string]string{}
	if body.PermissionName == "" {
		validationErrors["permission_name"] = "El nombre del permiso no puede estar vacío"
	} else if len(body.PermissionName) < 4 || len(body.PermissionName) > 50 {
		validationErrors["permission_name"] = "El nombre del permiso debe contener entre 4 y 50 caracteres"
	}

	if body.Description == "" {
		validationErrors["description"] = "La descripción no puede estar vacía"
	} else if len(body.Description) > 100 {
		validationErrors["description"] = "La descripción sólo puede contener hasta 100 caracteres"
	}

	if len(validationErrors) > 0 {
		return apperror.NewErrValidation(validationErrors)
	}

	permissionID, err := handler.permissionsService.Create(body.PermissionName, body.Description)
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusCreated, models.Permission{
		ID:          permissionID,
		Name:        body.PermissionName,
		Description: body.Description,
	})
}

func (handler *PermissionsHandler) GetPermissionByID(c *gin.Context) error {
	id := c.Param("id")

	permissionID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	permission := handler.permissionsService.GetPermissionByID(permissionID)
	if permission == nil {
		return apperror.NewErrPermissionNotFound()
	}

	return handler.JSONResponse(c, http.StatusOK, permission)
}

func (handler *PermissionsHandler) GetPermissionByName(c *gin.Context) error {
	permissionName := c.Param("permissionName")

	permission := handler.permissionsService.GetPermissionByName(permissionName)
	if permission == nil {
		return apperror.NewErrPermissionNotFound()
	}

	return handler.JSONResponse(c, http.StatusOK, permission)
}

func (handler *PermissionsHandler) GetPermissions(c *gin.Context) error {
	permissions := handler.permissionsService.GetPermissions()

	return handler.JSONResponse(c, http.StatusOK, permissions)
}

func (handler *PermissionsHandler) DeletePermission(c *gin.Context) error {
	permissionName := c.Param("permissionName")

	err := handler.permissionsService.DeletePermission(permissionName)
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusNoContent, nil)
}

func (handler *PermissionsHandler) GetPermissionsForUser(c *gin.Context) error {
	username := c.Param("username")

	user := handler.usersService.GetByUsername(username)
	if user == nil {
		return apperror.NewErrUserNotFound()
	}

	userPermissions := handler.permissionsService.GetPermissionsForUser(user.ID)

	permissionsNames := []string{}
	for _, userPermission := range userPermissions {
		permission := handler.permissionsService.GetPermissionByID(userPermission.PermissionID)
		if permission == nil {
			continue
		}

		permissionsNames = append(permissionsNames, permission.Name)
	}

	return handler.JSONResponse(c, http.StatusOK, permissionsNames)
}

func (handler *PermissionsHandler) GrantPermissionToUser(c *gin.Context) error {
	username := c.Param("username")
	permissionName := c.Param("permissionName")

	user := handler.usersService.GetByUsername(username)
	if user == nil {
		return apperror.NewErrUserNotFound()
	}

	err := handler.permissionsService.GrantPermissionToUser(user.ID, permissionName)
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusNoContent, nil)
}

func (handler *PermissionsHandler) RevokePermissionToUser(c *gin.Context) error {
	username := c.Param("username")
	permissionName := c.Param("permissionName")

	currentUser := c.MustGet("user").(models.User)
	if strings.EqualFold(currentUser.Username, username) {
		return apperror.NewErrCannotRevokeUserPermission()
	}

	user := handler.usersService.GetByUsername(username)
	if user == nil {
		return apperror.NewErrUserNotFound()
	}

	err := handler.permissionsService.RevokePermissionToUser(user.ID, permissionName)
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusNoContent, nil)
}

func NewPermissionsHandler(
	logger logger.Logger,

	permissionsService services.PermissionsService,
	usersService services.UsersService,
) *PermissionsHandler {
	return &PermissionsHandler{
		BaseHandler: BaseHandler{
			logger: logger,
		},

		permissionsService: permissionsService,
		usersService:       usersService,
	}
}
