package main

import (
	"rest-api/config"

	v1 "rest-api/handler/v1"
	"rest-api/middleware"

	"rest-api/repo"
	"rest-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB            = config.SetupDatabaseConnection()
	userRepo    repo.UserRepository = repo.NewUserRepo(db)
	authService service.AuthService = service.NewAuthService(userRepo)
	jwtService  service.JWTService  = service.NewJWTService()
	userService service.UserService = service.NewUserService(userRepo)
	authHandler v1.AuthHandler      = v1.NewAuthHandler(authService, jwtService, userService)
	userHandler v1.UserHandler      = v1.NewUserHandler(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	server := gin.Default()

	authRoutes := server.Group("/api/v1/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
	}
	userRoutes := server.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userHandler.Profile)
		userRoutes.PUT("/profile", userHandler.Update)
	}

	server.Run(":8081")
}
