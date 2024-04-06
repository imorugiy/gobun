package schema

import (
	"context"
	"database/sql"
)

type Model interface {
	ScanRows(ctx context.Context, rows *sql.Rows) (int, error)
	Value() interface{}
}
