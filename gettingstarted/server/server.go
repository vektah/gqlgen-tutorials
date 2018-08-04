package main

import (
	"log"
	"net/http"
	"os"

	gettingstarted "github.com/vektah/gqlgen-tutorials/gettingstarted"
	"github.com/vektah/gqlgen/handler"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(gettingstarted.NewExecutableSchema(gettingstarted.Config{Resolvers: &gettingstarted.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
