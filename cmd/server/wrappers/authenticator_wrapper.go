package wrappers

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/platform/authenticator"
	"go-crud-gin/internal/platform/logger"
	"go-crud-gin/internal/services"

	"github.com/gin-gonic/gin"
)

type authenticationHeaders struct {
	Authorization *string `header:"Authorization"`
}

type AuthenticatorWrapper struct {
	logger        logger.Logger
	authenticator authenticator.Authenticator
	usersService  services.UsersService
}

func (wrapper *AuthenticatorWrapper) Wrap(handler func(c *gin.Context) error, permissions []string) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		var headers authenticationHeaders

		err := c.ShouldBindHeader(&headers)
		if err == nil && headers.Authorization != nil {
			jwt, err := wrapper.authenticator.Authenticate(*headers.Authorization, permissions)
			if err != nil {
				return err
			}

			user := wrapper.usersService.GetByID(jwt.UserID)
			if user == nil {
				return apperror.NewErrUnauthorized()
			}

			c.Set("user", *user)
		} else if len(permissions) > 0 {
			return apperror.NewErrUnauthorized()
		}

		return handler(c)
	}
}

func NewAuthentiatorWrapper(
	logger logger.Logger,
	authenticator authenticator.Authenticator,
	usersService services.UsersService,
) *AuthenticatorWrapper {
	return &AuthenticatorWrapper{
		logger:        logger,
		authenticator: authenticator,
		usersService:  usersService,
	}
}
