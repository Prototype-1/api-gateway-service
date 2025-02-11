package main

import (
	"log"
	"github.com/Prototype-1/api-gateway-service/config" 
	"github.com/gin-gonic/gin"
	rt "github.com/Prototype-1/api-gateway-service/internal/router"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbl9pZCI6NiwiZXhwIjoxNzM5MTgyMTkyLCJyb2xlIjoiYWRtaW4ifQ.T5mQxR5rAdolLWhbDvZ6iPCb1o8JR6vOzqj1ps4lIps"
	secretKey := "RsaGPCEv4Oy1rRzVq8rT9vyPV29DZ4yQ4X_xd0QRKbM"

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println("Token validation error:", err)
		return
	}
	fmt.Println("Token valid:", token.Valid)
	
	config.LoadConfig()

	router := gin.Default()
	adminClient, userClient, routeClient, err := rt.InitGRPCClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	rt.SetupRoutes(router, adminClient, userClient, routeClient)

	log.Println("API Gateway running on port 8080")
	log.Fatal(router.Run(":8080"))
}
