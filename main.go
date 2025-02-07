package main

import (
	"log"
	"net/http"

	"github.com/Prototype-1/api-gateway-service/config" 
)

func main() {
	config.LoadConfig()
	log.Println("API Gateway running on port 8080")
	http.ListenAndServe(":8080", nil)
}
