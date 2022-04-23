package v1

import (
	"net/http"
	"strconv"

	"rest-api/common/obj"
	"rest-api/common/response"
	"rest-api/dto"
	"rest-api/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	jwtService  service.JWTService
	userService service.UserService
}

func NewAuthHandler(
	authService service.AuthService,
	jwtService service.JWTService,
	userService service.UserService,
) AuthHandler {
	return &authHandler{
		authService: authService,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (c *authHandler) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	err := ctx.ShouldBind(&loginRequest)

	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.authService.VerifyCredential(loginRequest.Email, loginRequest.Password)
	if err != nil {
		response := response.BuildErrorResponse("Failed to login", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user, _ := c.userService.FindUserByEmail(loginRequest.Email)

	token := c.jwtService.GenerateToken(strconv.FormatInt(user.ID, 10))
	user.Token = token
	response := response.BuildResponse(true, "OK!", user)
	ctx.JSON(http.StatusOK, response)

}

func (c *authHandler) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	err := ctx.ShouldBind(&registerRequest)
	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.CreateUser(registerRequest)
	if err != nil {
		response := response.BuildErrorResponse(err.Error(), err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	token := c.jwtService.GenerateToken(strconv.FormatInt(user.ID, 10))
	user.Token = token
	response := response.BuildResponse(true, "OK!", user)
	ctx.JSON(http.StatusCreated, response)

}

// package v1

// import (
// 	"rest-api/service"

// 	"github.com/labstack/echo/v4"
// )

// type AuthHandler interface {
// 	Login(ctx echo.Context) error
// 	Register(ctx echo.Context) error
// }

// type authHandler struct {
// 	userService service.UserService
// }

// func NewAuthHandler(userService service.UserService) AuthHandler {
// 	return &authHandler{
// 		userService: userService,
// 	}
// }

// func (a *authHandler) Login(ctx echo.Context) error {
// 	return nil
// }

// func (a *authHandler) Register(ctx echo.Context) error {
// 	return nil
// }
