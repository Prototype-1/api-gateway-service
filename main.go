package main

import (
	"log"
	"github.com/Prototype-1/api-gateway-service/config" 
	"github.com/gin-gonic/gin"
	"github.com/Prototype-1/api-gateway-service/internal/middleware"
)

func main() {
	config.LoadConfig()

	router := gin.Default()


	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware("admin"))
	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin Dashboard Accessed"})
		})
	}

	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware("user"))
	{
		userRoutes.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "User Profile Accessed"})
		})
	}

	log.Fatal(router.Run(":8080"))
}
