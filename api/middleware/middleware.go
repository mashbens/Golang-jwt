package middleware

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"

	"rest-api/common/obj"
	"rest-api/common/response"
	"rest-api/service"

	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(jwtService service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			header := c.Request().Header.Get("Authorization")
			token := jwtService.ValidateToken(header, c)

			// fmt.Println("header: ", header)
			// fmt.Println("token: ", token)
			if header == "" {
				response := response.BuildErrorResponse("Error", "Failed to get token", obj.EmptyObj{})
				c.JSON(http.StatusUnauthorized, response)
				return nil
			}

			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				log.Println("Claim[user_id]: ", claims["user_id"])
				log.Println("Claim[isuser]: ", claims["isuser"])

			} else {
				response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
				c.JSON(http.StatusUnauthorized, response)
				return nil
			}
			return next(c)
		}
	}
}
