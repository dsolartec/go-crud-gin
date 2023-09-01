package apperror

import (
	"fmt"
	"net/http"
)

const (
	ErrInternalServerCode    = "internal_server_error"
	ErrInternalServerMessage = "Error interno"

	ErrValidationCode    = "validation_error"
	ErrValdiationMessage = "Ha ocurrido un error validando la información"

	ErrUnauthorizedCode    = "unauthorized_user"
	ErrUnauthorizedMessage = "No tienes permitido consumir esta url"

	// Users
	ErrUserWrongAuthenticationCode    = "wrong_authentication"
	ErrUserWrongAuthenticationMessage = "El usuario o la contraseña no son correctos"

	ErrUserAlreadyExistsCode    = "user_already_exists"
	ErrUserAlreadyExistsMessage = "El nombre de usuario ya está en uso"

	ErrUserNotFoundCode    = "user_not_found"
	ErrUserNotFoundMessage = "El usuario no existe"

	ErrUserNotDeletableCode    = "cannot_delete_user"
	ErrUserNotDeletableMessage = "No puedes eliminar el usuario con el que estás autenticado"

	// Permissions
	ErrPermissionAlreadyExistsCode    = "permission_already_exists"
	ErrPermissionAlreadyExistsMessage = "El nombre del permiso ya está en uso"

	ErrPermissionNotFoundCode    = "permission_not_found"
	ErrPermissionNotFoundMessage = "El permiso no existe"

	ErrPermissionNotDeletableCode    = "permission_not_deletable"
	ErrPermissionNotDeletableMessage = "No puedes eliminar este permiso"

	// User permissions
	ErrUserAlreadyHasPermissionCode    = "user_has_permission"
	ErrUserAlreadyHasPermissionMessage = "El usuario ya posee este permiso"

	ErrUserPermissionNotFoundCode    = "user_not_has_permission"
	ErrUserPermissionNotFoundMessage = "El usuario no posee este permiso"

	ErrCannotRevokeUserPermissionCode    = "cannot_revoke_permission"
	ErrCannotRevokeUserPermissionMessage = "No puedes eliminarle un permiso al usuario con el que estás autenticado"
)

type AppError struct {
	StatusCode        int               `json:"status_code"`
	Code              string            `json:"code"`
	Message           string            `json:"message"`
	Details           *error            `json:"details,omitempty"`
	ValidationDetails map[string]string `json:"validation_details,omitempty"`
}

func (appError *AppError) Error() string {
	if appError.Details != nil {
		return fmt.Sprintf("%s - Error: %v", appError.Message, appError.Details)
	}

	return appError.Message
}

func NewErrInternalServerError(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       ErrInternalServerCode,
		Message:    ErrInternalServerMessage,
		Details:    &err,
	}
}

func NewErrValidation(validationDetails map[string]string) *AppError {
	return &AppError{
		StatusCode:        http.StatusBadRequest,
		Code:              ErrValidationCode,
		Message:           ErrValdiationMessage,
		ValidationDetails: validationDetails,
	}
}

func NewErrUnauthorized() *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       ErrUnauthorizedCode,
		Message:    ErrUnauthorizedMessage,
	}
}

// Users
func NewErrUserWrongAuthentication() *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       ErrUserWrongAuthenticationCode,
		Message:    ErrUserWrongAuthenticationMessage,
	}
}

func NewErrUserAlreadyExists() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrUserAlreadyExistsCode,
		Message:    ErrUserAlreadyExistsMessage,
	}
}

func NewErrUserNotFound() *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       ErrUserNotFoundCode,
		Message:    ErrUserNotFoundMessage,
	}
}

func NewErrUserNotDeletable() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrUserNotDeletableCode,
		Message:    ErrUserNotDeletableMessage,
	}
}

// Permissions
func NewErrPermissionAlreadyExists() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrPermissionAlreadyExistsCode,
		Message:    ErrPermissionAlreadyExistsMessage,
	}
}

func NewErrPermissionNotFound() *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       ErrPermissionNotFoundCode,
		Message:    ErrPermissionNotFoundMessage,
	}
}

func NewErrPermissionNotDeletable() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrPermissionNotDeletableCode,
		Message:    ErrPermissionNotDeletableMessage,
	}
}

// User permissions
func NewErrUserAlreadyHasPermission() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrUserAlreadyHasPermissionCode,
		Message:    ErrUserAlreadyHasPermissionMessage,
	}
}

func NewErrUserPermissionNotFound() *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       ErrUserPermissionNotFoundCode,
		Message:    ErrUserPermissionNotFoundMessage,
	}
}

func NewErrCannotRevokeUserPermission() *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       ErrCannotRevokeUserPermissionCode,
		Message:    ErrCannotRevokeUserPermissionMessage,
	}
}
