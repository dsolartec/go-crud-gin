package wrappers

import (
	"errors"
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

type ErrorWrapper struct {
	logger logger.Logger
}

func (wrapper *ErrorWrapper) Wrap(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err == nil {
			return
		}

		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			appErr = apperror.NewErrInternalServerError(err)
		}

		c.JSONP(appErr.StatusCode, appErr)
	}
}

func NewErrorWrapper(
	logger logger.Logger,
) *ErrorWrapper {
	return &ErrorWrapper{
		logger: logger,
	}
}
