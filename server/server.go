package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/ridhamtarpara/go-graphql-demo"
	"github.com/ridhamtarpara/go-graphql-demo/api/dal"
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
	dal.DBConn = db

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(go_graphql_demo.NewExecutableSchema(go_graphql_demo.Config{Resolvers: &go_graphql_demo.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}