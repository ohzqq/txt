package txt

import (
	"errors"
)

type Field struct {
	Tokens Tokens
}

var (
	NoMatchErr = errors.New(`no matches found`)
)

func NewField(toks Tokens) *Field {
	return &Field{Tokens: toks}
}

func (f *Field) Find(q string) (Tokens, error) {
	return f.Tokens.Find(q)
}

func (f *Field) Search(q string) (Tokens, error) {
	return f.Tokens.Search(q)
}

func (f *Field) FindByLabel(label string) (*Token, error) {
	return f.Tokens.FindByLabel(label)
}

func (f *Field) FindByValue(val string) (*Token, error) {
	return f.Tokens.FindByValue(val)
}

func (f *Field) FindByIndex(ti ...int) (Tokens, error) {
	return f.Tokens.FindByIndex(ti)
}

func (f *Field) Sort(order ...string) Tokens {
	o := "asc"
	if len(order) > 0 {
		o = order[0]
	}
	return f.Tokens.Sort(o)
}

func (f *Field) Asc() Tokens {
	return f.Tokens.Asc()
}

func (f *Field) Desc() Tokens {
	return f.Tokens.Desc()
}

func (toks *Field) String(i int) string {
	return toks.Tokens.String(i)
}

func (toks *Field) Len() int {
	return toks.Tokens.Len()
}
