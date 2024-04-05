package pgdialect

import (
	"database/sql"

	"github.com/imorugiy/go-project/dialect"
	"github.com/imorugiy/go-project/schema"
)

type Dialect struct {
	schema.BaseDialect
}

func New() *Dialect {
	d := new(Dialect)
	return d
}

func (d *Dialect) Init(*sql.DB) {}

func (d *Dialect) Name() dialect.Name {
	return dialect.PG
}
