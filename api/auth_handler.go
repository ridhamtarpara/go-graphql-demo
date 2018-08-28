package api

import (
	"net/http"
	"context"
	"github.com/ridhamtarpara/go-graphql-jobs/dal/firebase"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbConn := firebase.Connect()
		header := r.Header.Get("Authorization")
		// if auth is not available then proceed to resolver
		if header == "" {
			next.ServeHTTP(w, r)
		} else {
			fbData, err := dbConn.Auth.VerifyIDToken(dbConn.Context, header)
			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				// merge userID into request context
				ctx := context.WithValue(r.Context(), "userID", fbData.UID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
