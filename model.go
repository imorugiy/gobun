package main

import (
	"errors"
	"reflect"

	"github.com/imorugiy/go-project/schema"
)

var errNilModel = errors.New("bun: Model(nil)")

type Model = schema.Model

func newSingleModel(db *DB, dest interface{}) (Model, error) {
	return _newModel(db, dest, false)
}

func _newModel(db *DB, dest interface{}, scan bool) (Model, error) {
	switch dest := dest.(type) {
	case nil:
		return nil, errNilModel
	case Model:
		return dest, nil
	}

	v := reflect.ValueOf(dest)

	if v.IsNil() {
		typ := v.Type().Elem()
		if typ.Kind() == reflect.Struct {
			// TODO: here
		}
	}

	return nil, errors.New("Unhandled error")
}
