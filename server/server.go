package main

import (
	"database/sql"
	"github.com/99designs/gqlgen/handler"
	"github.com/ridhamtarpara/go-graphql-demo"
	"github.com/ridhamtarpara/go-graphql-demo/api/auth"
	"github.com/ridhamtarpara/go-graphql-demo/api/dal"
	"github.com/ridhamtarpara/go-graphql-demo/api/dataloaders"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := dal.Connect()
	checkErr(err)

	initDB(db)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	rootHandler:= dataloaders.DataloaderMiddleware(
		db,
		handler.GraphQL(
			go_graphql_demo.NewExecutableSchema(go_graphql_demo.NewRootResolvers(db)),
			handler.ComplexityLimit(200)),
	)
	http.Handle("/query", auth.AuthMiddleware(rootHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB(db *sql.DB) {
	dal.MustExec(db,"DROP TABLE IF EXISTS reviews")
	dal.MustExec(db,"DROP TABLE IF EXISTS screenshots")
	dal.MustExec(db,"DROP TABLE IF EXISTS videos")
	dal.MustExec(db,"DROP TABLE IF EXISTS users")
	dal.MustExec(db,"CREATE TABLE public.users (id SERIAL PRIMARY KEY, name varchar(255), email varchar(255))")
	dal.MustExec(db,"CREATE TABLE public.videos (id SERIAL PRIMARY KEY, name varchar(255), description varchar(255), url text,created_at TIMESTAMP, user_id int, FOREIGN KEY (user_id) REFERENCES users (id))")
	dal.MustExec(db,"CREATE TABLE public.screenshots (id SERIAL PRIMARY KEY, video_id int, url text, FOREIGN KEY (video_id) REFERENCES videos (id))")
	dal.MustExec(db,"CREATE TABLE public.reviews (id SERIAL PRIMARY KEY, video_id int,user_id int, description varchar(255), rating varchar(255), created_at TIMESTAMP, FOREIGN KEY (user_id) REFERENCES users (id), FOREIGN KEY (video_id) REFERENCES videos (id))")
	dal.MustExec(db,"INSERT INTO users(name, email) VALUES('Ridham', 'contact@ridham.me')")
	dal.MustExec(db,"INSERT INTO users(name, email) VALUES('Tushar', 'tushar@ridham.me')")
	dal.MustExec(db,"INSERT INTO users(name, email) VALUES('Dipen', 'dipen@ridham.me')")
	dal.MustExec(db,"INSERT INTO users(name, email) VALUES('Harsh', 'harsh@ridham.me')")
	dal.MustExec(db,"INSERT INTO users(name, email) VALUES('Priyank', 'priyank@ridham.me')")
}
