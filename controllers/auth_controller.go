package controllers

import (
	"fmt"
	"net/http"
	"start-speek/db"
	"start-speek/models"
	"start-speek/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"email"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(context *gin.Context) {
	var input RegisterInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	db.DB.Where("email = ?", input.Email).First(&existingUser)

	if existingUser.ID != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	db.DB.Create(&user)

	context.JSON(http.StatusOK, gin.H{"message": "User registered."})
}

func Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	result := db.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password 1",
		})
		return
	}

	valid := utils.CheckPassword(input.Password, user.Password)

	if !valid {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password 2",
		})
		return
	}

	accessToken, err := utils.GenerateToken(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token."})
		return
	}

	sessionKey := fmt.Sprintf("session:user:%d", user.ID)

	err = db.RedisClient.Set(
		db.Ctx,
		sessionKey,
		"active",
		7*24*time.Hour,
	).Err()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	refreshKey := fmt.Sprintf("refresh:user:%d", user.ID)

	err = db.RedisClient.Set(
		db.Ctx,
		refreshKey,
		refreshToken,
		7*24*time.Hour,
	).Err()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store refresh token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetMe(context *gin.Context) {

	userId, exists := context.Get("user_id")

	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User

	db.DB.First(&user, userId)

	if user.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func UpdateMe(context *gin.Context) {
	userId := context.GetUint("user_id")

	var input UpdateUserInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	db.DB.First(&user, userId)

	if user.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	db.DB.Save(&user)

	context.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func ChangePassword(context *gin.Context) {
	userId := context.GetUint("user_id")

	var input ChangePasswordInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	db.DB.First(&user, userId)

	if user.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	user.Password = hashedPassword

	db.DB.Save(&user)

	context.JSON(http.StatusOK, gin.H{"message": "Password change successfully."})
}

func LogOut(context *gin.Context) {
	userId := context.GetUint("user_id")

	sessionKey := fmt.Sprintf("session:user:%d", userId)

	err := db.RedisClient.Del(db.Ctx, sessionKey).Err()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout."})
		return
	}

	refreshKey := fmt.Sprintf("refresh:user:%d", userId)

	err = db.RedisClient.Get(db.Ctx, refreshKey).Err()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
}

func RefreshToken(context *gin.Context) {
	var input RefreshInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return utils.JWTKey, nil
	})

	if err != nil || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	// check redis
	refreshKey := fmt.Sprintf("refresh:user:%d", userId)

	storedToken, err := db.RedisClient.Get(db.Ctx, refreshKey).Result()

	if err != nil || storedToken != input.RefreshToken {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired."})
		return
	}

	// generate new access token
	newAccessToken, err := utils.GenerateRefreshToken(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
