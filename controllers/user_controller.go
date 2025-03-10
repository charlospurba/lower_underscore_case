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

// User struct untuk model data pengguna
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
}

// Data dummy untuk contoh
var users = []User{
	{ID: 1, Username: "john_doe", Email: "john@example.com", FirstName: "John", LastName: "Doe"},
	{ID: 2, Username: "jane_doe", Email: "jane@example.com", FirstName: "Jane", LastName: "Doe"},
}

// GetUsers mengembalikan daftar semua pengguna
// @Summary Get all users
// @Description Ambil semua pengguna yang tersedia
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /api/users [get]
// GetUsers mengembalikan daftar semua pengguna dari database
// GetUsers mengambil semua pengguna dari database
func GetUsers(c *gin.Context) {
	var users []models.User

	// Ambil data dari database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := database.DB.QueryContext(ctx, "SELECT id, username, email, first_name, last_name, age FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users", "details": err.Error()})
		return
	}
	defer rows.Close()

	// Loop untuk membaca hasil query
	for rows.Next() {
		var user models.User
		var age sql.NullInt32 // Variabel sementara untuk menangani NULL di kolom age

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning users", "details": err.Error()})
			return
		}

		// Jika age tidak NULL, simpan nilainya, jika NULL, biarkan user.Age tetap nil
		if age.Valid {
			ageInt := int(age.Int32)
			user.Age = &ageInt
		} else {
			user.Age = nil
		}

		users = append(users, user)
	}

	// Cek jika tidak ada user ditemukan
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No users found", "users": []models.User{}})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID mengambil pengguna berdasarkan ID
// @Summary Get user by ID
// @Description Ambil satu pengguna berdasarkan ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Ambil data berdasarkan ID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := database.DB.QueryRowContext(ctx, "SELECT id, username, email, first_name, last_name, age FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Age)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser menambahkan pengguna baru
// @Summary Create new user
// @Description Tambahkan pengguna baru ke dalam daftar
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User Data"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var newUser models.User

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	newUser.Password = hashedPassword

	// Simpan user ke database, termasuk Age
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO users (username, email, first_name, last_name, password, age) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = database.DB.QueryRowContext(ctx, query, newUser.Username, newUser.Email, newUser.FirstName, newUser.LastName, newUser.Password, newUser.Age).Scan(&newUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": newUser})
}

// DeleteUser menghapus pengguna berdasarkan ID
// @Summary Delete user
// @Description Hapus pengguna berdasarkan ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Hapus user dari database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := database.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Cek apakah user berhasil dihapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUser memperbarui data pengguna berdasarkan ID
// @Summary Update user by ID
// @Description Perbarui informasi pengguna berdasarkan ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body User true "Updated User Data"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser models.User

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data user di database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE users SET username = $1, email = $2, first_name = $3, last_name = $4, age = $5 WHERE id = $6"
	result, err := database.DB.ExecContext(ctx, query, updatedUser.Username, updatedUser.Email, updatedUser.FirstName, updatedUser.LastName, updatedUser.Age, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Cek apakah user berhasil diperbarui
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
