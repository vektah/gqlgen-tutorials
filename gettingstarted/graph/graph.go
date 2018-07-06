//go:generate gqlgen

package graph

import (
	"context"
	"fmt"
	"math/rand"
)

type User struct {
	ID   string
	Name string
}

type Todo struct {
	ID     string
	Text   string
	Done   bool
	UserID string
}

type MyApp struct {
	todos []Todo
}

func (a *MyApp) Query_todos(ctx context.Context) ([]Todo, error) {
	return a.todos, nil
}

func (a *MyApp) Mutation_createTodo(ctx context.Context, input NewTodo) (Todo, error) {
	todo := Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: input.User,
	}
	a.todos = append(a.todos, todo)
	return todo, nil
}

func (a *MyApp) Todo_user(ctx context.Context, it *Todo) (User, error) {
	return User{ID: it.UserID, Name: "user " + it.UserID}, nil
}
