package app

import (
	"go-crud-gin/cmd/server/handlers"
	"go-crud-gin/cmd/server/wrappers"
	authenticatorpkg "go-crud-gin/internal/platform/authenticator"
	loggerpkg "go-crud-gin/internal/platform/logger"
	"go-crud-gin/internal/services"

	"github.com/gin-gonic/gin"
)

type App interface {
	Run() error
}

type app struct {
	router        *gin.Engine
	authenticator authenticatorpkg.Authenticator
	logger        loggerpkg.Logger

	// Services
	usersService       services.UsersService
	permissionsService services.PermissionsService

	// Handlers
	authHandler        *handlers.AuthHandler
	usersHandler       *handlers.UsersHandler
	permissionsHandler *handlers.PermissionsHandler

	// Wrappers
	authenticatorWrapper *wrappers.AuthenticatorWrapper
	errorWrapper         *wrappers.ErrorWrapper
}

func (app *app) setupDependencies() {
	app.logger.Infof("[APP] Setting up dependencies...")

	// Handlers
	app.authHandler = handlers.NewAuthHandler(app.logger, app.authenticator, app.usersService, app.permissionsService)
	app.usersHandler = handlers.NewUsersHandler(app.logger, app.usersService)
	app.permissionsHandler = handlers.NewPermissionsHandler(app.logger, app.permissionsService, app.usersService)

	// Wrappers
	app.authenticatorWrapper = wrappers.NewAuthentiatorWrapper(app.logger, app.authenticator, app.usersService)
	app.errorWrapper = wrappers.NewErrorWrapper(app.logger)

	app.logger.Infof("[APP] Dependencies setted up!")
}

func (app *app) setupRouter() {
	app.logger.Infof("[APP] Setting up routes...")

	auth := app.router.Group("/auth")
	auth.POST("/logIn", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.authHandler.LogIn, []string{})))
	auth.POST("/signUp", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.authHandler.SignUp, []string{})))

	users := app.router.Group("/users")
	users.GET("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.usersHandler.GetUsers, []string{"users_read", "users_full"})))
	users.GET("/id/:id", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.usersHandler.GetUserByID, []string{"users_read", "users_full"})))

	userActions := users.Group("/username/:username")
	userActions.GET("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.usersHandler.GetUserByUsername, []string{"users_read", "users_full"})))
	userActions.DELETE("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.usersHandler.DeleteUser, []string{"users_write", "users_full"})))
	userActions.GET("/permissions", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.GetPermissionsForUser, []string{"users_read", "users_full"})))

	permissions := app.router.Group("/permissions")
	permissions.GET("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.GetPermissions, []string{"permissions_read", "permissions_full"})))
	permissions.POST("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.CreatePermission, []string{"permissions_write", "permissions_full"})))
	permissions.GET("/id/:id", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.GetPermissionByID, []string{"permissions_read", "permissions_full"})))
	permissions.GET("/name/:permissionName", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.GetPermissionByName, []string{"permissions_read", "permissions_full"})))
	permissions.DELETE("/name/:permissionName", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.DeletePermission, []string{"permissions_write", "permissions_full"})))

	userPermissions := userActions.Group("/permission/:permissionName")
	userPermissions.POST("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.GrantPermissionToUser, []string{"grant_permission"})))
	userPermissions.DELETE("/", app.errorWrapper.Wrap(app.authenticatorWrapper.Wrap(app.permissionsHandler.RevokePermissionToUser, []string{"revoke_permission"})))

	app.logger.Infof("[APP] Routes setted up!")
}

func (app *app) setup() {
	app.logger.Infof("[APP] Setting up application...")

	app.setupDependencies()
	app.setupRouter()

	app.logger.Infof("[APP] Application setted up!")
}

func (app *app) Run() error {
	return app.router.Run(":8080")
}

func newApp(
	router *gin.Engine,
	logger loggerpkg.Logger,
	authenticator authenticatorpkg.Authenticator,
) App {
	if router == nil {
		router = gin.Default()
	}

	if logger == nil {
		logger = loggerpkg.NewLocalLogger()
	}

	if authenticator == nil {
		authenticator = authenticatorpkg.NewLocalAuthenticator(logger)
	}

	// Services
	usersService := services.NewUsersService(logger)
	permissionsService := services.NewPermissionsService(logger)

	app := &app{
		router:        router,
		authenticator: authenticator,
		logger:        logger,

		// Services
		usersService:       usersService,
		permissionsService: permissionsService,
	}

	app.setup()

	return app
}
