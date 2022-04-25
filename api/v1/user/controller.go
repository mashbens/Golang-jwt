package user

import (
	"fmt"
	"net/http"
	"rest-api/common/response"
	"rest-api/dto"
	"rest-api/service"
	"strconv"

	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(
	userService service.UserService,
	jwtService service.JWTService,
) *UserController {
	return &UserController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (controller *UserController) getUserIDByHeader(c echo.Context) string {
	header := c.Request().Header.Get("Authorization")

	if header == "" {
		return fmt.Sprint("Error", "Failed to validate token")
	}
	token := controller.jwtService.ValidateToken(header, c)
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}

func (controller *UserController) Update(c echo.Context) error {
	var updateUserRequest dto.UpdateUserRequest

	err := c.Bind(&updateUserRequest)
	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", "Invalid request body", nil)
		return c.JSON(http.StatusBadRequest, response)
	}
	id := controller.getUserIDByHeader(c)

	if id == "" {
		response := response.BuildErrorResponse("Failed to process request", "Invalid user id", nil)
		return c.JSON(http.StatusBadRequest, response)
	}
	_id, _ := strconv.ParseInt(id, 10, 64)
	updateUserRequest.ID = _id
	res, err := controller.userService.UpdateUser(updateUserRequest)

	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := response.BuildResponse(true, "User updated", res)
	return c.JSON(http.StatusOK, response)
}

func (controller *UserController) Profile(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	token := controller.jwtService.ValidateToken(header, c)
	if header == "" {
		response := response.BuildErrorResponse("Failed to process request", "Failed to validate token", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user, err := controller.userService.FindUserByID(id)

	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := response.BuildResponse(true, "User found", user)
	return c.JSON(http.StatusOK, response)
}
