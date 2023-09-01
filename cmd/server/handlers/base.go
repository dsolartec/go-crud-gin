package handlers

import (
	"go-crud-gin/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	logger logger.Logger
}

func (handler *BaseHandler) JSONResponse(c *gin.Context, statusCode int, data any) error {
	c.JSONP(statusCode, data)
	return nil
}
