package api

import (
	"rest-api/api/middleware"
	"rest-api/api/v1/auth"
	"rest-api/api/v1/user"

	"rest-api/service"

	"github.com/labstack/echo/v4"
)

var jwtService service.JWTService = service.NewJWTService()

type Router struct {
	Auth *auth.AuthController
	User *user.UserController
}

func RegisterRoutes(e *echo.Echo, controller *Router) {

	authRoutes := e.Group("/api/v1/auth")
	authRoutes.POST("/login", controller.Auth.Login)
	authRoutes.POST("/register", controller.Auth.Register)

	userRoutes := e.Group("/api/v1/user", middleware.AuthorizeJWT(jwtService))
	// userRoutes.Use(middleware.AuthorizeJWT(service.NewJWTService()))
	userRoutes.GET("/profile", controller.User.Profile)
	userRoutes.PUT("/profile", controller.User.Update)

	// userRoutes := e.Group("/api/v1/user", middleware.AuthorizeJWT(service.NewJWTService()))
	// // userRoutes.Use(middleware.AuthorizeJWT(service.NewJWTService()))
	// userRoutes.GET("/profile", controller.User.Profile)
	// userRoutes.PUT("/profile", controller.User.Update)

}
