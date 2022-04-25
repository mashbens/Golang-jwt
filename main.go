package main

import (
	api "rest-api/api"
	"rest-api/config"

	// v1 "rest-api/handler/v1"

	"rest-api/repo"
	"rest-api/service"

	router "rest-api/api"
	// "rest-api/api/middleware"
	"rest-api/api/v1/auth"
	"rest-api/api/v1/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB            = config.SetupDatabaseConnection()
	userRepo    repo.UserRepository = repo.NewUserRepo(db)
	authService service.AuthService = service.NewAuthService(userRepo)
	jwtService  service.JWTService  = service.NewJWTService()
	userService service.UserService = service.NewUserService(userRepo)
	// authHandler auth.AuthController = *auth.NewAuthController(authService, jwtService, userService)
	// userHandler user.UserController = *user.NewUserController(userService, jwtService)
)

// type Controller struct {
// 	Auth *auth.AuthController
// 	User *user.UserController
// }

// func RegisterRoutes(e *echo.Echo, controller *Controller) {

// 	authRoutes := e.Group("/api/v1/auth")
// 	authRoutes.POST("/login", controller.Auth.Login)
// 	authRoutes.POST("/register", controller.Auth.Register)

// 	userRoutes := e.Group("/api/v1/user", middleware.AuthorizeJWT(jwtService))
// 	userRoutes.Use(middleware.AuthorizeJWT(service.NewJWTService()))
// 	userRoutes.GET("/profile", controller.User.Profile)
// 	userRoutes.PUT("/profile", controller.User.Update)

// }

// func main() {
// 	defer config.CloseDatabaseConnection(db)

// 	e := echo.New()
// 	controller := &Controller{
// 		Auth: auth.NewAuthController(authService, jwtService, userService),
// 		User: user.NewUserController(userService, jwtService),
// 	}
// 	RegisterRoutes(e, controller)

// 	e.Logger.Fatal(e.Start(":8081"))
// }

func main() {
	defer config.CloseDatabaseConnection(db)

	e := echo.New()
	controler := router.Router{
		Auth: auth.NewAuthController(authService, jwtService, userService),
		User: user.NewUserController(userService, jwtService),
	}
	api.RegisterRoutes(e, &controler)
	e.Logger.Fatal(e.Start(":8081"))
}
