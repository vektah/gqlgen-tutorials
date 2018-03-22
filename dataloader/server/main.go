package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/vektah/gqlgen-tutorials/dataloader"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/")
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

	queryHandler := handler.GraphQL(dataloader.MakeExecutableSchema(dataloader.New(db)))

	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", dataloader.DataloaderMiddleware(db, queryHandler))

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mustExec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
