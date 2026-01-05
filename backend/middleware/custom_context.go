package middleware

import (
	"expenses-tracker/models"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func CustomContextMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			cc := &models.CustomContext{
				Context:  c,
				UserID:   0,
				Email:    "",
				UserName: "",
			}

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header required"})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid authorization header format"})
			}

			tokenString := parts[1]
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key-change-in-production"
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			// Verify user still exists and get user info
			var user models.M_user
			if err := db.First(&user, claims.UserID).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
			}

			cc.UserID = claims.UserID
			cc.Email = claims.Email
			cc.UserName = user.Name

			return next(cc)
		}
	}
}

func GetCustomContext(c echo.Context) *models.CustomContext {
	if cc, ok := c.(*models.CustomContext); ok {
		return cc
	}
	return &models.CustomContext{
		Context:  c,
		UserID:   0,
		Email:    "",
		UserName: "",
	}
}
