package main

type baseQuery struct{}

func (q *baseQuery) setModel(modeli interface{}) {

}

type whereBaseQuery struct {
	baseQuery
}
