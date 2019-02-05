package auth

import (
	"context"
	"github.com/ridhamtarpara/go-graphql-demo"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			// Write your fancy token introspection logic here and if valid user then pass appropriate key in header
			// IMPORTANT: DO NOT HANDLE UNAUTHORISED USER HERE
			ctx = context.WithValue(ctx,go_graphql_demo.UserIDCtxKey, auth)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}