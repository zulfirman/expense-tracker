package middleware

import (
	"net/http"
	"os"
	"strings"

	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CustomContextMiddleware(userRepo *repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid token"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key-change-in-production"
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.ErrUnauthorized
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			userIDFloat, ok := claims["userId"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID in token"})
			}

			userID := uint(userIDFloat)
			user, err := userRepo.GetByID(userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
			}

			cc := &model.CustomContext{
				Context: c,
				UserID:  user.ID,
				Email:   user.Email,
				UserName: user.Name,
			}

			return next(cc)
		}
	}
}

// Helper to retrieve our custom context
func GetCustomContext(c echo.Context) *model.CustomContext {
	if cc, ok := c.(*model.CustomContext); ok {
		return cc
	}
	return &model.CustomContext{Context: c}
}


