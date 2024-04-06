package schema

import "github.com/imorugiy/go-project/dialect"

type Formatter struct {
	dialect Dialect
	args    *namedArgList
}

func NewFormatter(dialect Dialect) Formatter {
	return Formatter{dialect: dialect}
}

func (f Formatter) IsNop() bool {
	return f.dialect.Name() == dialect.Invalid
}

func (f Formatter) WithArg(arg NamedArgAppender) Formatter {
	return Formatter{dialect: f.dialect, args: f.args.WithArg(arg)}
}

type NamedArgAppender interface {
	AppendNamedArg(fmter Formatter, b []byte, name string) ([]byte, bool)
}

type namedArgList struct {
	arg  NamedArgAppender
	next *namedArgList
}

func (l *namedArgList) WithArg(arg NamedArgAppender) *namedArgList {
	return &namedArgList{
		arg:  arg,
		next: l,
	}
}

func (l *namedArgList) AppendNamedArg(fmter Formatter, b []byte, name string) ([]byte, bool) {
	for l != nil && l.arg != nil {
		if b, ok := l.arg.AppendNamedArg(fmter, b, name); ok {
			return b, true
		}
		l = l.next
	}
	return b, false
}
