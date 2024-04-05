package main

import "context"

type SelectQuery struct {
	whereBaseQuery
}

func NewSelectQuery(db *DB) *SelectQuery {
	return &SelectQuery{}
}

func (q *SelectQuery) Model(model interface{}) *SelectQuery {
	q.setModel(model)
	return q
}

func (q *SelectQuery) Count(ctx context.Context) (int, error) {
	return 1, nil
}
