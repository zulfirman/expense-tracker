package handler

import (
	"crypto/rand"
	"encoding/base64"
	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	db               *gorm.DB
}

func NewAuthHandler(userRepo *repository.UserRepository, refreshTokenRepo *repository.RefreshTokenRepository, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		db:               db,
	}
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
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         model.M_user `json:"user"`
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
	_, err := h.userRepo.GetByEmail(req.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already registered"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}

	// Create user
	user := model.M_user{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.userRepo.Create(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	// Seed default categories for the new user
	if err := repository.SeedDefaultCategories(h.db, user.ID); err != nil {
		// Log error but don't fail signup - categories can be added later
		// In production, you might want to log this to a monitoring service
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
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
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
		User:         *user,
	})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateCurrencyRequest struct {
	Currency string `json:"currency" validate:"required,oneof=IDR USD EUR JPY"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Check if email is being changed and if it's already taken
	if req.Email != user.Email {
		exists, err := h.userRepo.EmailExists(req.Email, userID)
		if err == nil && exists {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Email already in use"})
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := h.userRepo.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateCurrency(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req UpdateCurrencyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	user.Currency = req.Currency

	if err := h.userRepo.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update currency"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Current password is incorrect"})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}

	user.Password = string(hashedPassword)

	if err := h.userRepo.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
}

func (h *AuthHandler) generateAccessToken(userID uint, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Minute * 3).Unix(), // 3 minutes
		"type":   "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (h *AuthHandler) generateRefreshToken(userID uint) (string, error) {
	// Generate a random token string
	tokenString := generateRandomString(32)

	// Store refresh token in database
	refreshToken := model.M_refresh_token{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
	}

	if err := h.refreshTokenRepo.Create(&refreshToken); err != nil {
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
	refreshToken, err := h.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
	}

	// Get user
	user, err := h.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate access token"})
	}

	// Optionally rotate refresh token (delete old, create new)
	h.refreshTokenRepo.Delete(refreshToken)
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
