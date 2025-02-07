package router

import (
	"context"
	"log"
	"net/http"
	"github.com/Prototype-1/api-gateway-service/internal/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	uPB "github.com/Prototype-1/api-gateway-service/proto/user"
	aPB "github.com/Prototype-1/api-gateway-service/proto/admin"
	rPB "github.com/Prototype-1/api-gateway-service/proto/routes"
)

func SetupRouter(mux *runtime.ServeMux) error {
	userConn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user-auth-service: %v", err)
	}

	adminConn, err := grpc.DialContext(context.Background(), "localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to admin-auth-service: %v", err)
	}

	routeConn, err := grpc.DialContext(context.Background(), "localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to admin-routes-service: %v", err)
	}

	if err := uPB.RegisterUserServiceHandler(context.Background(), mux, userConn); err != nil {
		return err
	}

	if err := aPB.RegisterAdminServiceHandler(context.Background(), mux, adminConn); err != nil {
		return err
	}

	if err := rPB.RegisterRouteServiceHandler(context.Background(), mux, routeConn); err != nil {
		return err
	}

	return nil
}


func MiddlewareHandler(mux *runtime.ServeMux) http.Handler {
	apiRouter := http.NewServeMux()
	apiRouter.Handle("/", mux)

	// Protect specific routes
	apiRouter.Handle("/admin/", middleware.AuthMiddleware(mux, "admin"))
	apiRouter.Handle("/user/", middleware.AuthMiddleware(mux, "user"))

	return apiRouter
}
