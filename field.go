package txt

type Field struct {
	Tokens Tokens
}

func NewField(toks Tokens) *Field {
	return &Field{Tokens: toks}
}

func (toks *Field) Find(q string) Tokens {
	return toks.Tokens.Find(q)
}

func (toks *Field) FuzzyFind(q string) Tokens {
	return toks.Tokens.FuzzyFind(q)
}

func (toks *Field) String(i int) string {
	return toks.Tokens.String(i)
}

func (toks *Field) Len() int {
	return toks.Tokens.Len()
}
