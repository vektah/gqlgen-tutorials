package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vektah/gqlgen-tutorials/gettingstarted/graph"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	app := &graph.MyApp{}
	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", handler.GraphQL(graph.MakeExecutableSchema(app)))

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
