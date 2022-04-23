package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"rest-api/common/obj"
	"rest-api/common/response"
	"rest-api/dto"
	"rest-api/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserHandler interface {
	Profile(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserHandler(
	userService service.UserService,
	jwtService service.JWTService,
) UserHandler {
	return &userHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userHandler) getUserIDByHeader(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	token := c.jwtService.ValidateToken(header, ctx)

	if token == nil {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}

func (c *userHandler) Update(ctx *gin.Context) {
	var updateUserRequest dto.UpdateUserRequest

	err := ctx.ShouldBind(&updateUserRequest)
	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	id := c.getUserIDByHeader(ctx)

	if id == "" {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	_id, _ := strconv.ParseInt(id, 0, 64)
	updateUserRequest.ID = _id
	res, err := c.userService.UpdateUser(updateUserRequest)

	if err != nil {
		response := response.BuildErrorResponse("Error", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	response := response.BuildResponse(true, "OK", res)
	ctx.JSON(http.StatusOK, response)

}

func (c *userHandler) Profile(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	token := c.jwtService.ValidateToken(header, ctx)

	if token == nil {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user, err := c.userService.FindUserByID(id)

	if err != nil {
		response := response.BuildErrorResponse("Error", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	res := response.BuildResponse(true, "OK", user)
	ctx.JSON(http.StatusOK, res)
}

// package v1

// import (
// 	"rest-api/service"

// 	"github.com/labstack/echo/v4"
// )

// type UserHandler interface {
// 	Update(ctx echo.Context) error
// 	Profile(ctx echo.Context) error
// 	// GetUserIdByHeader(header string) (int64, error)
// 	// Update(id int64, user *entity.User) (*entity.User, error)
// 	// Profile(id int64) (*entity.User, error)
// 	// Profile(ctx *gin.Context)
// 	// Update(ctx *gin.Context)
// }

// type userHandler struct {
// 	userService service.UserService
// }

// func NewUserHandler(userService service.UserService) UserHandler {
// 	return &userHandler{
// 		userService: userService,
// 	}
// }

// func (u *userHandler) Update(ctx echo.Context) error {
// 	return nil
// }

// func (u *userHandler) Profile(ctx echo.Context) error {
// 	return nil
// }
