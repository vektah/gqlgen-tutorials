//go:generate gqlgen -schema ./schema.graphql -typemap types.json
//go:generate dataloaden github.com/vektah/gqlgen-tutorials/dataloader.User

package dataloader

import (
	"context"
	"database/sql"

	"net/http"

	"time"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Resolver struct {
	db *sql.DB
}

const userLoaderKey = "userloader"

func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := UserLoader{
			maxBatch: 100,
			wait:     1 * time.Millisecond,
			fetch: func(ids []int) ([]*User, []error) {
				placeholders := make([]string, len(ids))
				args := make([]interface{}, len(ids))
				for i := 0; i < len(ids); i++ {
					placeholders[i] = "?"
					args[i] = i
				}

				res := logAndQuery(db,
					"SELECT id, name from dataloader_example.user WHERE id IN ("+
						strings.Join(placeholders, ",")+")",
					args...,
				)
				defer res.Close()

				users := make([]*User, len(ids))
				i := 0
				for res.Next() {
					users[i] = &User{}
					err := res.Scan(&users[i].ID, &users[i].Name)
					if err != nil {
						panic(err)
					}
					i++
				}

				return users, nil
			},
		}
		ctx := context.WithValue(r.Context(), userLoaderKey, &userloader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value(userLoaderKey).(*UserLoader)
}

func New(db *sql.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

type Todo struct {
	ID     string
	Todo   string
	UserID int
}

func (r *Resolver) Query_todos(ctx context.Context) ([]Todo, error) {
	res := logAndQuery(r.db, "SELECT id, todo, user_id FROM dataloader_example.todo")
	defer res.Close()

	var todos []Todo
	for res.Next() {
		var todo Todo
		if err := res.Scan(&todo.ID, &todo.Todo, &todo.UserID); err != nil {
			panic(err)
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *Resolver) Todo_userRaw(ctx context.Context, obj *Todo) (*User, error) {
	res := logAndQuery(r.db, "SELECT id, name FROM dataloader_example.user WHERE id = ?", obj.UserID)
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}
	var user User
	if err := res.Scan(&user.ID, &user.Name); err != nil {
		panic(err)
	}
	return &user, nil
}

func (r *Resolver) Todo_userLoader(ctx context.Context, obj *Todo) (*User, error) {
	return getUserLoader(ctx).Load(obj.UserID)
}
