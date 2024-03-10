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

func (f *Field) FuzzyFind(q string) (Tokens, error) {
	return f.Tokens.FuzzyFind(q)
}

func (f *Field) FindByLabel(label string) (*Token, error) {
	return f.Tokens.FindByLabel(label)
}

func (f *Field) FindByValue(val string) (*Token, error) {
	return f.Tokens.FindByValue(val)
}

func (toks *Field) String(i int) string {
	return toks.Tokens.String(i)
}

func (toks *Field) Len() int {
	return toks.Tokens.Len()
}
