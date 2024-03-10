package txt

import (
	"fmt"
	"strings"

	"github.com/sahilm/fuzzy"
)

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

func Tokenize(ana *Analyzer, text string) (Tokens, error) {
	var (
		toks   Tokens
		tokens []string
	)

	if text == "" {
		return toks, EmptyStrErr
	}

	if ana.sep == nil {
		tokens = []string{text}
	} else {
		tokens = strings.FieldsFunc(text, ana.sep)
	}

	if len(tokens) == 0 {
		return toks, FieldsFuncErr
	}

	for _, label := range tokens {
		tok := label
		for _, norm := range ana.normalizers {
			tok = norm(tok)
		}

		if ana.RmStopWords() {
			if ana.IsStopWord(strings.ToLower(tok)) {
				continue
			}
		}

		if tok != "" {
			toks = append(toks, NewToken(label, tok))
		}
	}

	return toks, nil
}

func newMatch(str string, idx int) fuzzy.Match {
	return fuzzy.Match{
		Str:   str,
		Index: idx,
	}
}

func (toks Tokens) Find(q string) (Tokens, error) {
	var tokens Tokens
	for i, tok := range toks {
		if tok.Value == q {
			tok.Match = newMatch(tok.Label, i)
			tokens = append(tokens, tok)
		}
	}
	if tokens.Len() > 0 {
		return tokens, nil
	}

	return nil, fmt.Errorf("%w for query '%s'\n", NoMatchErr, q)
}

func (toks Tokens) FuzzyFind(q string) (Tokens, error) {
	var tokens Tokens
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

func (toks Tokens) FindByLabel(label string) (*Token, error) {
	for i, token := range toks {
		if token.Label == label {
			token.Match = newMatch(token.Label, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for label '%s'\n", NoMatchErr, label)
}

func (toks Tokens) FindByValue(val string) (*Token, error) {
	for i, token := range toks {
		if token.Value == val {
			token.Match = newMatch(token.Value, i)
			return token, nil
		}
	}
	return nil, fmt.Errorf("%w for val '%s'\n", NoMatchErr, val)
}

func (toks Tokens) String(i int) string {
	return toks[i].Value
}

func (toks Tokens) Len() int {
	return len(toks)
}
