package main

import (
	"context"
	"fmt"

	"github.com/imorugiy/go-project/schema"
)

type SelectQuery struct {
	whereBaseQuery
}

func NewSelectQuery(db *DB) *SelectQuery {
	return &SelectQuery{
		whereBaseQuery: whereBaseQuery{
			baseQuery: baseQuery{
				db: db,
			},
		},
	}
}

func (q *SelectQuery) Model(model interface{}) *SelectQuery {
	q.setModel(model)
	return q
}

func (q *SelectQuery) Count(ctx context.Context) (int, error) {
	if q.err != nil {
		return 0, q.err
	}

	qq := countQuery{q}

	queryBytes, err := qq.AppendQuery(q.db.fmter, nil)
	if err != nil {
		return 0, nil
	}

	fmt.Println(queryBytes, err)
	// query := internal.String(queryBytes)
	// ctx, event := q.db.beforeQuery(ctx, qq, query, nil,)
	return 1, nil
}

func (q *SelectQuery) appendQuery(fmter schema.Formatter, b []byte, count bool) (_ []byte, err error) {
	if q.err != nil {
		return nil, q.err
	}

	fmter = formatterWithModel(fmter, q)

	if count {
		b = append(b, "count(*)"...)
	}

	return b, nil
}

type countQuery struct {
	*SelectQuery
}

func (q countQuery) AppendQuery(fmter schema.Formatter, b []byte) (_ []byte, err error) {
	if q.err != nil {
		return nil, q.err
	}
	return q.appendQuery(fmter, b, true)
}
