package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// CustomContextMiddleware verifies JWT access tokens and transparently handles
// refresh using a database-backed refresh token store. When an access token is
// expired but a valid refresh token is provided, it will:
//   - validate the refresh token from X-Refresh-Token
//   - enforce a usage limit using UsedCount
//   - generate a new access and refresh token pair
//   - persist the new refresh token and increment usage
//   - send new tokens via X-Token and X-Refresh-Token headers
//   - continue to the next handler without requiring the client to retry
func CustomContextMiddleware(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			refreshHeader := c.Request().Header.Get("X-Refresh-Token")
			deviceID := c.Request().Header.Get("X-Device-Id")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid token"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key-change-in-production"
			}

			// Helper to build custom context and continue request
			setUserContextAndNext := func(user *model.M_user) error {
				// Read workspace ID from header
				workspaceIDStr := c.Request().Header.Get("X-Workspace-Id")
				var workspaceID uint = 0
				if workspaceIDStr != "" {
					if id, err := strconv.ParseUint(workspaceIDStr, 10, 64); err == nil {
						workspaceID = uint(id)
					}
				}
				
				cc := &model.CustomContext{
					Context:     c,
					UserID:      user.ID,
					Email:       user.Email,
					UserName:    user.Name,
					WorkspaceID: workspaceID,
				}
				return next(cc)
			}

			// Parse access token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.ErrUnauthorized
				}
				return []byte(secret), nil
			})

			if err == nil && token.Valid {
				// Valid access token
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
				}

				if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token type"})
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

				return setUserContextAndNext(user)
			}

			// Check if token is expired and try refresh flow
			if ve, ok := err.(*jwt.ValidationError); ok && (ve.Errors&jwt.ValidationErrorExpired) != 0 {
				// Access token expired
				if refreshHeader == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "authorization token expired and refresh token missing"})
				}

				// Look up refresh token in DB
				rt, err := refreshRepo.GetByToken(refreshHeader)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid or expired refresh token"})
				}

				if time.Now().After(rt.ExpiresAt) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "refresh token expired"})
				}

				// Enforce usage limit (e.g., 30 times)
				if rt.UsedCount >= 30 {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "refresh token exceeded usage limit"})
				}

				// Load user
				user, err := userRepo.GetByID(rt.UserID)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "user not found for refresh token"})
				}

				// Generate new access token
				claims := jwt.MapClaims{
					"userId": user.ID,
					"email":  user.Email,
					"exp":    time.Now().Add(time.Minute * 3).Unix(),
					"type":   "access",
				}
				accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				newAccessToken, err := accessToken.SignedString([]byte(secret))
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "failed to generate new access token"})
				}

				// Optionally rotate refresh token: here we keep the same token string,
				// but increment usage counter and update device ID if provided.
				rt.UsedCount++
				if deviceID != "" {
					rt.DeviceID = deviceID
				}
				if err := refreshRepo.Save(rt); err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "failed to update refresh token usage"})
				}

				// Send new tokens via headers (no need for client to retry request)
				c.Response().Header().Set("X-Token", newAccessToken)
				c.Response().Header().Set("X-Refresh-Token", rt.Token)

				// Continue with request using the user from refresh token
				return setUserContextAndNext(user)
			}

			// Other invalid cases
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
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
