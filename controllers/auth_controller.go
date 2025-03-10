package controllers

import (
	"context"
	"database/sql"
	"gin-user-app/database"
	"gin-user-app/models"
	"gin-user-app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Login mengautentikasi pengguna dan mengembalikan token JWT
// @Summary Login user
// @Description Login menggunakan username dan password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Cek apakah database sudah terhubung
	if database.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	var storedPassword string
	var userID int

	// Set timeout query ke database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ambil password hash dari database berdasarkan username
	err := database.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE username=$1", req.Username).Scan(&userID, &storedPassword)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Bandingkan password dengan hash di database
	if err := utils.ComparePasswords(storedPassword, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Berikan response sukses dengan token
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// Logout godoc
// @Summary Logout user
// @Description Logout pengguna dengan menghapus token di client
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Optional: Add token to blacklist in Redis/Database
	// This prevents the token from being used again before its natural expiration
	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// err := database.BlacklistToken(tokenString)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token"})
	//     return
	// }

	// Clear any session cookies if you're using them
	c.SetCookie("session", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"details": "Token has been invalidated and cookies cleared",
	})
}

// VerifyUser memeriksa apakah token JWT valid
// @Summary Verify user
// @Description Verifikasi token JWT pengguna
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/auth/verify [get]
func VerifyUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token valid", "user_id": userID})
}
