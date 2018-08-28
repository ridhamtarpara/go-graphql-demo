package dal

import (
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
)

func GetID(DBConn *db.Client, Context context.Context, path string) (string, error) {
	ref, err := DBConn.NewRef(path).Push(Context, nil)
	return ref.Key, err
}
