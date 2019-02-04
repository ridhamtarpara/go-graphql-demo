package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "docker"
	DB_NAME     = "postgres"
)

var once sync.Once

func Connect() (*sql.DB, error){
	var db *sql.DB
	var err error
	once.Do(func() {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME)
		db, err = sql.Open("postgres", dbinfo)
	})
	return db, err
}

// Connect with CockroachDB and open DB connection.
var DBConn *sql.DB