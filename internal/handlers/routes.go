package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/Prototype-1/api-gateway-service/proto"
	"github.com/Prototype-1/api-gateway-service/internal/middleware"
)


func ProxyRequest(client pb.AdminServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
	
		if c.FullPath() == "/admin/login" {
			var loginReq pb.LoginRequest
			if err := c.ShouldBindJSON(&loginReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Call gRPC service method
			resp, err := client.Login(context.Background(), &loginReq)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with gRPC call"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token":   resp.GetToken(),
				"success": resp.GetSuccess(),
			})
		}
	}
}

func SetupRoutes(router *gin.Engine, client pb.AdminServiceClient) {
	router.POST("/admin/signup", ProxyRequest(client))
	router.POST("/admin/login", ProxyRequest(client)) 

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware()) 
	{
		protected.GET("/users", ProxyRequest(client))
		protected.POST("/users/block", ProxyRequest(client))
		protected.POST("/users/unblock", ProxyRequest(client))
	}
}
