package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/ridhamtarpara/go-graphql-jobs/graph"
	"github.com/ridhamtarpara/go-graphql-jobs/api"
)

const defaultPort = "8090"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	queryHandler :=http.HandlerFunc(handler.Playground("Job", "/query"))
	rootHandler := http.HandlerFunc(handler.GraphQL(graph.NewExecutableSchema(api.NewRootResolvers())))

	http.Handle("/", queryHandler)
	http.Handle("/query", api.AuthHandler(rootHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
