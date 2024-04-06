package main

import (
	"database/sql"

	"github.com/imorugiy/go-project/schema"
)

type DB struct {
	*sql.DB

	fmter schema.Formatter

	dialect schema.Dialect
}

func NewDB(sqldb *sql.DB, dialect schema.Dialect) *DB {
	dialect.Init(sqldb)

	db := &DB{
		DB:      sqldb,
		dialect: dialect,
		fmter:   schema.NewFormatter(dialect),
	}

	return db
}

func (db *DB) NewSelect() *SelectQuery {
	return NewSelectQuery(db)
}
