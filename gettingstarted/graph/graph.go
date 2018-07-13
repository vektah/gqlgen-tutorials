//go:generate gqlgen -typemap types.json -schema ../schema.graphql

package graph

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/vektah/gqlgen-tutorials/gettingstarted/model"
)

type App struct {
	todos []model.Todo
}

func (a *App) Mutation() MutationResolver {
	return &mutationResolver{a}
}

func (a *App) Query() QueryResolver {
	return &queryResolver{a}
}

func (a *App) Todo() TodoResolver {
	return &todoResolver{a}
}

type queryResolver struct{ *App }

func (a *queryResolver) Todos(ctx context.Context) ([]model.Todo, error) {
	return a.todos, nil
}

type mutationResolver struct{ *App }

func (a *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (model.Todo, error) {
	todo := model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: input.User,
	}
	a.todos = append(a.todos, todo)
	return todo, nil
}

type todoResolver struct{ *App }

func (a *todoResolver) User(ctx context.Context, it *model.Todo) (model.User, error) {
	return model.User{ID: it.UserID, Name: "user " + it.UserID}, nil
}
