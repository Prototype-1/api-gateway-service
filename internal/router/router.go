package router

import (
	"context"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	adminpb "github.com/Prototype-1/api-gateway-service/proto/admin"
	userpb "github.com/Prototype-1/api-gateway-service/proto/user"
	routespb "github.com/Prototype-1/api-gateway-service/proto/routes"
	"github.com/Prototype-1/api-gateway-service/internal/middleware"
)

func InitGRPCClients() (adminpb.AdminServiceClient, userpb.UserServiceClient, routespb.RouteServiceClient, error) {
	adminConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Admin Auth Service: %v", err)
		return nil, nil, nil, err
	}
	adminClient := adminpb.NewAdminServiceClient(adminConn)
	userConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Auth Service: %v", err)
		return nil, nil, nil, err
	}
	userClient := userpb.NewUserServiceClient(userConn)

	routesConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Routes Service: %v", err)
		return nil, nil, nil, err
	}
	routeClient := routespb.NewRouteServiceClient(routesConn)

	return adminClient, userClient, routeClient, nil
}

func SetupRoutes(router *gin.Engine, adminClient adminpb.AdminServiceClient, userClient userpb.UserServiceClient, routeClient routespb.RouteServiceClient) {
	router.POST("/user/signup", func(c *gin.Context) {
		var req userpb.SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := userClient.Signup(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.POST("/user/login", func(c *gin.Context) {
		var req userpb.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := userClient.Login(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.POST("/admin/signup", func(c *gin.Context) {
		var req adminpb.AdminSignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := adminClient.AdminSignup(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.POST("/admin/login", func(c *gin.Context) {
		var req adminpb.AdminLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := adminClient.AdminLogin(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware("admin"))
	
	adminRoutes.GET("/users", func(c *gin.Context) {
		resp, err := userClient.GetAllUsers(context.Background(), &userpb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	adminRoutes.POST("/users/block", func(c *gin.Context) {
		var req userpb.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := userClient.BlockUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	
	router.GET("/routes/get-all-routes", func(c *gin.Context) {
		resp, err := routeClient.GetAllRoutes(context.Background(), &routespb.GetAllRoutesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	
	adminRoutes.POST("/routes/add", middleware.AuthMiddleware("admin"), func(c *gin.Context) {
		var req routespb.AddRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := routeClient.AddRoute(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
}