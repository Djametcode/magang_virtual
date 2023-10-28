package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/main/databases"
	"example.com/main/helper"
	"example.com/main/models"
	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
)

func GetAllUser(ctx *gin.Context) {
	var user models.User

	databases.DB.Find(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func RegistUser(ctx *gin.Context) {
	user := models.User{}

	// Assuming you're trying to get the user data from the request JSON
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	hashedPass, err := helper.HashPassword(user.Password)

	if err != nil {
		fmt.Println("error hashing pass", err)
		return
	}

	user.Password = hashedPass;

	existingUser := models.User{}
	if err := databases.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email or Username already exists"})
		return
	}

	databases.DB.Create(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func LoginUser(ctx *gin.Context) {
	var user models.User

    if err := ctx.ShouldBind(&user); err != nil {
        ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	plainPass := user.Password

    if err := databases.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email not registered yet"})
        return
    }

    userExistPass := user.Password
	emailUser := user.Email
	userId := user.ID

    if !helper.CheckPassword(plainPass, userExistPass) {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
        return
    }

	token, err := helper.GenerateJWT(userId, emailUser)

	if err != nil {
		fmt.Println("jwt error")
	}

	type UserResponse struct {
        ID        uint      `json:"ID"`
        Username  string    `json:"username"`
        Email     string    `json:"email"`
        Photos    []string  `json:"Photos"`
        CreatedAt time.Time `json:"CreatedAt"`
        UpdatedAt time.Time `json:"UpdatedAt"`
    }

	userResponse := UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    ctx.JSON(http.StatusOK, gin.H{
        "user": userResponse,
		"token": token,
    })
}

func UpdateUser(ctx *gin.Context) {
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
        return
    }

    var user models.User

    // Retrieve user ID from the URL parameter
    userID := ctx.Param("userId")

    // TODO: Get the user's email and user ID from the JWT token
    // Example: email := ctx.GetString("email")
    //          userId := ctx.GetUint("userId")

    // Ensure you have user authentication and authorization in place
    // to ensure the user is allowed to perform the update.

    // Check if a file is attached
    file, header, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files attached"})
        return
    }

    // Process and save the file to Cloudinary
    // Upload the file to Cloudinary using the Cloudinary SDK
    result, err := cld.Upload.Upload(file, cloudinary.UploadParams{})
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to Cloudinary"})
        return
    }

    // Update the user's photo with the Cloudinary URL
    user.Photos = append(user.Photos, models.Photo{URL: result.SecureURL})

    // Save the user data back to the database
    if err := databases.DB.Save(&user).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "User photo updated successfully"})
}

