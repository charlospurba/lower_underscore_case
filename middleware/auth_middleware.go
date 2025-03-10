package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"gin-user-app/utils" // Sesuaikan dengan module kamu
	"github.com/gin-gonic/gin"
)

// AuthMiddleware untuk verifikasi token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debug: cek request URL dulu
		fmt.Println("DEBUG: Incoming request for:", c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("ERROR: Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Cek format Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			fmt.Println("ERROR: Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Verifikasi token
		token := tokenParts[1]
		claims, err := utils.VerifyToken(token)
		if err != nil {
			fmt.Println("ERROR: Invalid or expired token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Simpan user_id di context
		c.Set("user_id", claims["user_id"])
		fmt.Println("User ID from token:", claims["user_id"])

		// Lanjutkan request
		c.Next()
	}
}
