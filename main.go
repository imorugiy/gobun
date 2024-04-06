package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imorugiy/go-project/dialect/pgdialect"
	"github.com/imorugiy/go-project/driver/pgdriver"
)

type User struct{}

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector())
	db := NewDB(sqldb, pgdialect.New())

	r, err := db.NewSelect().Model((*User)(nil)).Count(context.Background())
	fmt.Println(r, err)
	// err := db.Ping()
	// fmt.Println(err)

}
