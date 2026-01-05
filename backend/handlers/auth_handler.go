package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"expenses-tracker/middleware"
	"expenses-tracker/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type SignupRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
	User         models.M_user `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Check if user already exists
	var existingUser models.M_user
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already registered"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}

	// Create user
	user := models.M_user{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	// Generate access token (3 minutes) and refresh token (7 days)
	accessToken, err := h.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate access token"})
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Find user
	var user models.M_user
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}

	// Generate access token (3 minutes) and refresh token (7 days)
	accessToken, err := h.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate access token"})
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var user models.M_user
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var user models.M_user
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Check if email is being changed and if it's already taken
	if req.Email != user.Email {
		var existingUser models.M_user
		if err := h.db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Email already in use"})
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := h.db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) generateAccessToken(userID uint, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
		"type":   "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (h *AuthHandler) generateRefreshToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	// Generate a random token string
	tokenString := generateRandomString(32)

	// Store refresh token in database
	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
	}

	if err := h.db.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Find refresh token in database
	var refreshToken models.RefreshToken
	if err := h.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
	}

	// Get user
	var user models.M_user
	if err := h.db.First(&user, refreshToken.UserID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate access token"})
	}

	// Optionally rotate refresh token (delete old, create new)
	h.db.Unscoped().Delete(&refreshToken)
	newRefreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
	}

	return c.JSON(http.StatusOK, RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	})
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback: use timestamp + random bytes as string
		return time.Now().Format("20060102150405") + string(b)
	}
	encoded := base64.URLEncoding.EncodeToString(b)
	if len(encoded) >= length {
		return encoded[:length]
	}
	return encoded
}
