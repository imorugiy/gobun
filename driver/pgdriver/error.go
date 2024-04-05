package pgdriver

import "fmt"

type Error struct {
	m map[byte]string
}

func (err Error) Field(k byte) string {
	return err.m[k]
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %s (SQLSTATE=%s)", err.Field('S'), err.Field('M'), err.Field('C'))
}
