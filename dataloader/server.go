package main

import (
	"database/sql"
	"fmt"
	"gqlgen-tutorials/dataloader/dataloader"
	"gqlgen-tutorials/dataloader/graph"
	"gqlgen-tutorials/dataloader/graph/generated"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
)

const defaultPort = "8080"

func main() {
	db, err := sql.Open("mysql", "root@tcp(mysql.dockervm:49926)/")
	if err != nil {
		panic(err)
	}

	db.Exec("DROP DATABASE dataloader_example")
	mustExec(db, "CREATE DATABASE dataloader_example")
	mustExec(db, "CREATE TABLE dataloader_example.user (id int NOT NULL AUTO_INCREMENT, name varchar(255), PRIMARY KEY(id))")
	mustExec(db, "CREATE TABLE dataloader_example.todo (id int NOT NULL AUTO_INCREMENT, todo varchar(255), user_id int, PRIMARY KEY(id))")

	for i := 0; i < 5; i++ {
		mustExec(db, "INSERT INTO dataloader_example.user (name) VALUES (?)", i)
	}

	for i := 0; i < 20; i++ {
		mustExec(db, "INSERT INTO dataloader_example.todo (todo, user_id) VALUES (?, ?)", fmt.Sprintf("Todo %d", i), (i+1)%5)
	}

	db.Exec("INSERT INTO todo (todo, user_id)")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Conn: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", dataloader.Middleware(db, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustExec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
