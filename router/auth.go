package router

import (
	"net/http"

	"example.com/main/controllers"
	"example.com/main/middlewares"
	"example.com/main/models"
	"github.com/gin-gonic/gin"
)

func AuthRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/users/get-all-user", controllers.GetAllUser)
	r.POST("/users/register", controllers.RegistUser)
	r.POST("/users/login", controllers.LoginUser)
	r.POST("/test", middlewares.AuthMiddleware(), func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, user.Email)
	})
	r.PUT("/users/:userId", middlewares.AuthMiddleware(), controllers.UpdateUser)
	return r
}