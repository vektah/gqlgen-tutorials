//go:generate go run github.com/vektah/dataloaden UserLoader int *gqlgen-tutorials/dataloader/graph/model.User

package dataloader

import (
	"context"
	"database/sql"
	"gqlgen-tutorials/dataloader/db"
	"gqlgen-tutorials/dataloader/graph/model"
	"net/http"
	"strings"
	"time"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById UserLoader
}

func Middleware(conn *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([]*model.User, []error) {
					placeholders := make([]string, len(ids))
					args := make([]interface{}, len(ids))
					for i := 0; i < len(ids); i++ {
						placeholders[i] = "?"
						args[i] = i
					}

					res := db.LogAndQuery(conn,
						"SELECT id, name from dataloader_example.user WHERE id IN ("+strings.Join(placeholders, ",")+")",
						args...,
					)
					defer res.Close()

					userById := map[int]*model.User{}
					for res.Next() {
						user := model.User{}
						err := res.Scan(&user.ID, &user.Name)
						if err != nil {
							panic(err)
						}
						userById[user.ID] = &user
					}

					users := make([]*model.User, len(ids))
					for i, id := range ids {
						users[i] = userById[id]
						i++
					}

					return users, nil
				},
			},
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
