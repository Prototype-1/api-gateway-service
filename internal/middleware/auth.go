package middleware

import (
	"context"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
	cf "github.com/Prototype-1/api-gateway-service/config"
)

type UserClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", http.ErrNoCookie
	}

	return parts[1], nil
}

func VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return [] byte (cf.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func AuthMiddleware(next http.Handler, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := ExtractToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Role != requiredRole {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func AttachMetadata(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}
