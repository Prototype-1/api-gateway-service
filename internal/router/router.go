package router

import (
	"context"
	"log"
	"net/http"
	"strings"
	"github.com/Prototype-1/api-gateway-service/internal/middleware"
	adminpb "github.com/Prototype-1/api-gateway-service/proto/admin"
	routespb "github.com/Prototype-1/api-gateway-service/proto/routes"
	userpb "github.com/Prototype-1/api-gateway-service/proto/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func InitGRPCClients() (adminpb.AdminServiceClient, userpb.UserServiceClient, routespb.RouteServiceClient, error) {
	adminConn, err := grpc.Dial("admin-auth-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Admin Auth Service: %v", err)
		return nil, nil, nil, err
	}
	adminClient := adminpb.NewAdminServiceClient(adminConn)
	userConn, err := grpc.Dial("user-auth-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Auth Service: %v", err)
		return nil, nil, nil, err
	}
	userClient := userpb.NewUserServiceClient(userConn)

	routesConn, err := grpc.Dial("admin-routes-service:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	
	router.GET("/routes/get-all-routes", func(c *gin.Context) {
		resp, err := routeClient.GetAllRoutes(context.Background(), &routespb.GetAllRoutesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	
	adminRoutes.POST("/routes/add", func(c *gin.Context) {
		var req routespb.AddRouteRequest
	
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			return
		}
	
		token = strings.TrimPrefix(token, "Bearer ")
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", token)
	
		resp, err := routeClient.AddRoute(ctx, &req)
		if err != nil {
			log.Printf("gRPC AddRoute call failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, resp)
	})
	
	adminRoutes.PUT("/routes/update", func(c *gin.Context) {
		var req routespb.UpdateRouteRequest
	
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			return
		}
	
		token = strings.TrimPrefix(token, "Bearer ")
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", token)
	
		resp, err := routeClient.UpdateRoute(ctx, &req)
		if err != nil {
			log.Printf("gRPC UpdateRoute call failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, resp)
	})

	adminRoutes.DELETE("/routes/delete", func(c *gin.Context) {
		var req routespb.DeleteRouteRequest
	
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			return
		}
	
		token = strings.TrimPrefix(token, "Bearer ")
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", token)
	
		resp, err := routeClient.DeleteRoute(ctx, &req)
		if err != nil {
			log.Printf("gRPC DeleteRoute call failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	
	adminRoutes.POST("/users/unblock", func(c *gin.Context) {
		var req userpb.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := userClient.UnblockUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	
	adminRoutes.POST("/users/suspend", func(c *gin.Context) {
		var req userpb.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		resp, err := userClient.SuspendUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC call failed"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
}
