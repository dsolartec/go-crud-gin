package handlers

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/models"
	"go-crud-gin/internal/platform/logger"
	"go-crud-gin/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	BaseHandler

	usersService services.UsersService
}

func (handler *UsersHandler) GetUsers(c *gin.Context) error {
	users := handler.usersService.GetUsers()

	return handler.JSONResponse(c, http.StatusOK, users)
}

func (handler *UsersHandler) GetUserByID(c *gin.Context) error {
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return err
	}

	user := handler.usersService.GetByID(userID)
	if user == nil {
		return apperror.NewErrUserNotFound()
	}

	return handler.JSONResponse(c, http.StatusOK, user)
}

func (handler *UsersHandler) GetUserByUsername(c *gin.Context) error {
	username := c.Param("username")

	user := handler.usersService.GetByUsername(username)
	if user == nil {
		return apperror.NewErrUserNotFound()
	}

	return handler.JSONResponse(c, http.StatusOK, user)
}

func (handler *UsersHandler) DeleteUser(c *gin.Context) error {
	username := c.Param("username")

	currentUser := c.MustGet("user").(models.User)
	if strings.EqualFold(currentUser.Username, username) {
		return apperror.NewErrUserNotDeletable()
	}

	err := handler.usersService.DeleteUser(username)
	if err != nil {
		return err
	}

	return handler.JSONResponse(c, http.StatusNoContent, nil)
}

func NewUsersHandler(
	logger logger.Logger,

	usersService services.UsersService,
) *UsersHandler {
	return &UsersHandler{
		BaseHandler: BaseHandler{
			logger: logger,
		},

		usersService: usersService,
	}
}
