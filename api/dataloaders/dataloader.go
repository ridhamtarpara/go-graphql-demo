package dataloaders

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/ridhamtarpara/go-graphql-demo/api"
	"github.com/ridhamtarpara/go-graphql-demo/api/dal"
	"github.com/ridhamtarpara/go-graphql-demo/api/errors"
	"log"
	"net/http"
	"time"
)

type ctxKeyType struct { name string }

var CtxKey = ctxKeyType{"dataloaderctx"}

type Loaders struct {
	UserByID  			*UserLoader
}

func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := UserLoader{
			wait : 1 * time.Millisecond,
			maxBatch: 100,
			fetch: func(ids []int) ([]*api.User, []error) {
				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, name, email from users WHERE id = ?"
				} else {
					sqlQuery = "SELECT id, name, email from users WHERE id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Println(err)
				}
				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := dal.LogAndQuery(db, sqlQuery, arguments...)
				defer rows.Close();
				if err != nil {
					log.Println(err)
				}
				userById := map[int]*api.User{}

				for rows.Next() {
					user:= api.User{}
					if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
						errors.DebugPrintf(err)
						return nil, []error{errors.InternalServerError}
					}
					userById[user.ID] = &user
				}

				users := make([]*api.User, len(ids))
				for i, id := range ids {
					users[i] = userById[id]
					i++
				}

				return users, nil
			},
		}
		ctx := context.WithValue(r.Context(), CtxKey, &userloader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

//func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		loader := Loaders{}
//
//		wait := 1 * time.Millisecond
//
//		loader.UserByID = &UserLoader{
//			wait : wait,
//			maxBatch: 100,
//			fetch: func(ids []int) ([]*api.User, []error) {
//				placeholders := make([]string, len(ids))
//				args := make([]interface{}, len(ids))
//				for i := 0; i < len(ids); i++ {
//					placeholders[i] = "?"
//					args[i] = i
//				}
//
//				res,err := dal.LogAndQuery(db, "SELECT id, name from user WHERE id IN ("+
//						strings.Join(placeholders, ",")+")",
//					args...,
//				)
//				defer res.Close()
//
//				if err!= nil {
//					return nil, nil
//				}
//
//				users := make([]*api.User, len(ids))
//				i := 0
//				for res.Next() {
//					users[i] = &api.User{}
//					err := res.Scan(&users[i].ID, &users[i].Name)
//					if err != nil {
//						panic(err)
//					}
//					i++
//				}
//
//				return users, nil
//			},
//		}
//		ctx = context.WithValue(ctx, CtxKey, loader)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func CtxLoaders(ctx context.Context) Loaders {
//	return ctx.Value(CtxKey).(Loaders)
//}