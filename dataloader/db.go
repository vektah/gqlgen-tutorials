package dataloader

import (
	"database/sql"
	"fmt"
)

func logAndQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {
	fmt.Println(query)
	res, err := db.Query(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}
