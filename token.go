package txt

import (
	"errors"
	"fmt"
	"slices"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	fuzzy.Match
}

type Tokens struct {
	toks []*Token
}

type AllTokens []*Token

var (
	FieldsFuncErr = errors.New("strings.FieldsFunc returned an empty slice or the string was empty")
	EmptyStrErr   = errors.New("empty string")
	NoMatchErr    = errors.New(`no matches found`)
)

func NewToken(label, val string) *Token {
	return &Token{
		Value: val,
		Label: label,
	}
}

func NewTokens() *Tokens {
	return &Tokens{}
}

func (t *Tokens) AddString(str string) *Tokens {
	t.AddToken(NewToken(str, str))
	return t
}

func (t *Tokens) SetStringSlice(str []string) *Tokens {
	//t.AddToken(NewToken(str, str))
	return t
}

func (t *Tokens) SetTokens(toks []*Token) *Tokens {
	t.toks = toks
	return t
}

func (t *Tokens) AddToken(tok *Token) *Tokens {
	t.toks = append(t.toks, tok)
	return t
}

func (toks AllTokens) Find(q string) (AllTokens, error) {
	var tokens AllTokens
	for i, tok := range toks {
		if tok.Value == q {
			tok.Match = newMatch(tok.Value, i)
			tokens = append(tokens, tok)
		} else if tok.Label == q {
			tok.Match = newMatch(tok.Label, i)
			tokens = append(tokens, tok)
		}
	}
	if tokens.Len() > 0 {
		return tokens, nil
	}

	return nil, fmt.Errorf("%w for query '%s'\n", NoMatchErr, q)
}

func (toks AllTokens) Without(sw AllTokens) AllTokens {
	return lo.Without(toks, sw...)
}

func (toks AllTokens) Search(q string) (AllTokens, error) {
	var tokens AllTokens
	for _, m := range fuzzy.FindFrom(q, toks) {
		tok := toks[m.Index]
		tok.Match = m
		tok.Match.Str = tok.Label
		tokens = append(tokens, tok)
	}

	if tokens.Len() > 0 {
		return tokens, nil
	}
	return nil, fmt.Errorf("%w for query '%s'\n", NoMatchErr, q)
}

func (toks AllTokens) FindByLabel(label string) (*Token, error) {
	for i, token := range toks {
		if token.Label == label {
			token.Match = newMatch(token.Label, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for label '%s'\n", NoMatchErr, label)
}

func (toks AllTokens) FindByValue(val string) (*Token, error) {
	for i, token := range toks {
		if token.Value == val {
			token.Match = newMatch(token.Value, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for val '%s'\n", NoMatchErr, val)
}

func (toks AllTokens) FindByIndex(ti []int) (AllTokens, error) {
	var tokens AllTokens
	for _, tok := range ti {
		if tok < toks.Len() {
			tokens = append(tokens, toks[tok])
		}
	}
	if tokens.Len() > 0 {
		return tokens, nil
	}
	return nil, fmt.Errorf("%w for indices %v\n", NoMatchErr, ti)
}

func (toks AllTokens) Values() []string {
	vals := make([]string, toks.Len())
	for i, tok := range toks {
		vals[i] = tok.Value
	}
	return vals
}

func (toks AllTokens) Labels() []string {
	vals := make([]string, toks.Len())
	for i, tok := range toks {
		vals[i] = tok.Label
	}
	return vals
}

func (toks AllTokens) String(i int) string {
	return toks[i].Value
}

func (toks AllTokens) Len() int {
	return len(toks)
}

func (toks AllTokens) Sort(cmp func(a, b *Token) int, order string) AllTokens {
	tokens := toks
	slices.SortStableFunc(tokens, cmp)
	if order == "desc" {
		slices.Reverse(tokens)
	}
	return tokens
}

func (toks AllTokens) SortStable(cmp func(a, b *Token) int, order string) AllTokens {
	tokens := toks
	slices.SortStableFunc(tokens, cmp)
	if order == "desc" {
		slices.Reverse(tokens)
	}
	return tokens
}

func (toks AllTokens) SortAlphaAsc() AllTokens {
	return toks.Sort(SortByAlphaFunc, "asc")
}

func (toks AllTokens) SortAlphaDesc() AllTokens {
	return toks.Sort(SortByAlphaFunc, "desc")
}

func SortByAlphaFunc(a *Token, b *Token) int {
	switch {
	case a.Value > b.Value:
		return 1
	case a.Value == b.Value:
		return 0
	default:
		return -1
	}
}

func newMatch(str string, idx int) fuzzy.Match {
	return fuzzy.Match{
		Str:   str,
		Index: idx,
	}
}
