package handlers

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/platform/authenticator"
	"go-crud-gin/internal/platform/logger"
	"go-crud-gin/internal/requests"
	"go-crud-gin/internal/responses"
	"go-crud-gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler

	authenticator      authenticator.Authenticator
	usersService       services.UsersService
	permissionsService services.PermissionsService
}

func (handler *AuthHandler) LogIn(c *gin.Context) error {
	var body *requests.LogInRequest
	if err := c.ShouldBind(&body); err != nil {
		return err
	}

	user := handler.usersService.GetByUsername(body.Username)
	if user == nil || user.Password != body.Password {
		return apperror.NewErrUserWrongAuthentication()
	}

	permissions := []string{}
	for _, userPermission := range handler.permissionsService.GetPermissionsForUser(user.ID) {
		permission := handler.permissionsService.GetPermissionByID(userPermission.PermissionID)
		if permission == nil {
			continue
		}

		permissions = append(permissions, permission.Name)
	}

	tokenStr, err := handler.authenticator.GetToken(authenticator.AuthenticatorToken{
		UserID:      user.ID,
		Permissions: permissions,
	})
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusOK, responses.LogInResponse{
		AccessToken: tokenStr,
	})
}

func (handler *AuthHandler) SignUp(c *gin.Context) error {
	var body *requests.SignUpRequest
	if err := c.ShouldBind(&body); err != nil {
		return err
	}

	validationErrors := map[string]string{}
	if body.Username == "" {
		validationErrors["username"] = "El nombre de usuario no puede estar vacío"
	} else if len(body.Username) < 4 || len(body.Username) > 15 {
		validationErrors["username"] = "El nombre de usuario debe contener entre 4 y 15 caracteres"
	}

	if body.Password == "" {
		validationErrors["password"] = "La contraseña no puede estar vacía"
	} else if len(body.Password) < 8 || len(body.Password) > 40 {
		validationErrors["password"] = "La contraseña debe contener entre 8 y 40 caracteres"
	}

	if len(validationErrors) > 0 {
		return apperror.NewErrValidation(validationErrors)
	}

	userID, err := handler.usersService.Create(body.Username, body.Password)
	if err != nil {
		return err
	}

	permissions := []string{}
	for _, userPermission := range handler.permissionsService.GetPermissionsForUser(userID) {
		permission := handler.permissionsService.GetPermissionByID(userPermission.PermissionID)
		if permission == nil {
			continue
		}

		permissions = append(permissions, permission.Name)
	}

	tokenStr, err := handler.authenticator.GetToken(authenticator.AuthenticatorToken{
		UserID:      userID,
		Permissions: permissions,
	})
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusCreated, responses.SignUpResponse{
		AccessToken: tokenStr,
		UserID:      userID,
	})
}

func NewAuthHandler(
	logger logger.Logger,

	authenticator authenticator.Authenticator,
	usersService services.UsersService,
	permissionsService services.PermissionsService,
) *AuthHandler {
	return &AuthHandler{
		BaseHandler: BaseHandler{
			logger: logger,
		},

		authenticator:      authenticator,
		usersService:       usersService,
		permissionsService: permissionsService,
	}
}
