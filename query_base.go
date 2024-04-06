package main

import (
	"fmt"

	"github.com/imorugiy/go-project/schema"
)

type baseQuery struct {
	db *DB

	model Model
	err   error
}

func (q *baseQuery) setModel(modeli interface{}) {
	model, err := newSingleModel(q.db, modeli)
	fmt.Println(model, err)
	if err != nil {
		q.setErr(err)
		return
	}

	q.model = model
}

func (q *baseQuery) setErr(err error) {
	if q.err == nil {
		q.err = err
	}
}

func (q *baseQuery) hasTables() bool {
	return false
	// return q.modelHasTableName() || len()
}

func (q *baseQuery) AppendNamedArg(fmter schema.Formatter, b []byte, name string) ([]byte, bool) {
	return b, false
}

type whereBaseQuery struct {
	baseQuery
}

func formatterWithModel(fmter schema.Formatter, model schema.NamedArgAppender) schema.Formatter {
	if fmter.IsNop() {
		return fmter
	}
	return fmter.WithArg(model)
}
