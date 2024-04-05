package schema

import (
	"database/sql"

	"github.com/imorugiy/go-project/dialect"
)

type Dialect interface {
	Init(db *sql.DB)

	Name() dialect.Name
}

type BaseDialect struct{}

type nopDialect struct {
	BaseDialect
}

func newNopDialect() *nopDialect {
	d := new(nopDialect)
	return d
}

func (d *nopDialect) Init(db *sql.DB) {}

func (d *nopDialect) Name() dialect.Name {
	return dialect.Invalid
}
