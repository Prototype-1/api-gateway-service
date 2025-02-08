package main

import (
	"log"
	"github.com/Prototype-1/api-gateway-service/config" 
	"github.com/gin-gonic/gin"
	rt "github.com/Prototype-1/api-gateway-service/internal/router"
)

func main() {
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
