package main

import (
	"example.com/main/databases"
	"example.com/main/router"
	"github.com/gin-gonic/gin"
)

func main() {
	databases.InitDatabase()
	r := gin.Default()

	r.Use(router.AuthRouter().HandleContext)
	r.Run(":8080")
}