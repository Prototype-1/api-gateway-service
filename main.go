package main

import (
	"fmt"
	"log"

	"github.com/Prototype-1/api-gateway-service/internal/handlers"
	pb "github.com/Prototype-1/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to admin service: %v", err)
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	router := gin.Default()


	handlers.SetupRoutes(router, client)

	serverPort := ":8080"
	fmt.Println("API Gateway running on port", serverPort)
	log.Fatal(router.Run(serverPort))
}
