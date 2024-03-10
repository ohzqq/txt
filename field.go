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

func (f *Field) AlphaAsc() Tokens {
	return f.Tokens.AlphaAsc()
}

func (f *Field) AlphaDesc() Tokens {
	return f.Tokens.AlphaDesc()
}

func (toks *Field) String(i int) string {
	return toks.Tokens.String(i)
}

func (toks *Field) Len() int {
	return toks.Tokens.Len()
}
