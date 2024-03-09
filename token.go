package txt

import "github.com/sahilm/fuzzy"

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	fuzzy.Match
}

type Tokens []*Token

func NewToken(label, val string) *Token {
	return &Token{
		Value: val,
		Label: label,
	}
}

func (toks Tokens) Find(q string) Tokens {
	var tokens Tokens
	for i, tok := range toks {
		if tok.Value == q {
			tok.Match = fuzzy.Match{
				Index: i,
				Str:   tok.Label,
			}
			tokens = append(tokens, tok)
		}
	}
	return tokens
}

func (toks Tokens) FuzzyFind(q string) Tokens {
	var tokens Tokens
	for _, m := range fuzzy.FindFrom(q, toks) {
		tok := toks[m.Index]
		tok.Match = m
		tok.Match.Str = tok.Label
		tokens = append(tokens, tok)
	}
	return tokens
}

func (toks Tokens) String(i int) string {
	return toks[i].Value
}

func (toks Tokens) Len() int {
	return len(toks)
}
