package main

import (
	"database/sql"
	"fmt"

	"github.com/imorugiy/go-project/dialect/pgdialect"
	"github.com/imorugiy/go-project/driver/pgdriver"
)

type User struct{}

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector())
	db := NewDB(sqldb, pgdialect.New())

	// r, err := db.NewSelect().Model((*User)(nil)).Count(context.Background())

	err := db.Ping()
	fmt.Println(err)

}
