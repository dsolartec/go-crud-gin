package app

import (
	authenticatorpkg "go-crud-gin/internal/platform/authenticator"
	loggerpkg "go-crud-gin/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

type AppBuilder interface {
	Build() App
	WithRouter(router *gin.Engine) *appBuilder
	WithLogger(logger loggerpkg.Logger) *appBuilder
	WithAuthenticator(authenticator authenticatorpkg.Authenticator) *appBuilder
}

type appBuilder struct {
	router        *gin.Engine
	authenticator authenticatorpkg.Authenticator
	logger        loggerpkg.Logger
}

func (builder *appBuilder) WithRouter(router *gin.Engine) *appBuilder {
	builder.router = router
	return builder
}

func (builder *appBuilder) WithLogger(logger loggerpkg.Logger) *appBuilder {
	builder.logger = logger
	return builder
}

func (builder *appBuilder) WithAuthenticator(authenticator authenticatorpkg.Authenticator) *appBuilder {
	builder.authenticator = authenticator
	return builder
}

func (builder *appBuilder) Build() App {
	return newApp(
		builder.router,
		builder.logger,
		builder.authenticator,
	)
}

func NewAppBuilder() AppBuilder {
	return &appBuilder{}
}
