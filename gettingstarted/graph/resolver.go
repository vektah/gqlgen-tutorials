//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package graph

import "github.com/vektah/gqlgen-todos/graph/model"

type Resolver struct {
	todos []*model.Todo
}
