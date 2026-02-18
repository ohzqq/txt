package txt

import (
	"slices"

	"github.com/ohzqq/txt/sep"
)

type Analyzer struct {
	stopWords   []string
	sep         sep.Func
	normalizers []Normalizer
}

func New(opts ...Option) *Analyzer {
	ana := &Analyzer{
		sep: sep.OnWhitespace,
	}

	for _, opt := range opts {
		opt(ana)
	}

	return ana
}

func Keywords() *Analyzer {
	return New(ToLower).Keywords()
}

func NewNormalizer(opts ...Option) *Analyzer {
	ana := New(
		ToLower,
		WithoutPunct,
	)
	for _, opt := range opts {
		opt(ana)
	}
	return ana
}

func (ana *Analyzer) WithNormalizer(normies ...Normalizer) *Analyzer {
	ana.normalizers = normies
	return ana
}

func (ana *Analyzer) AddNormalizer(normies ...Normalizer) *Analyzer {
	ana.normalizers = append(ana.normalizers, normies...)
	return ana
}

func (ana *Analyzer) Tokenize(text string) (AllTokens, error) {
	var (
		tokz   AllTokens
		tokens []string
		//toks   = NewTokens()
	)

	if text == "" {
		return tokz, EmptyStrErr
	}

	tokens = Split(text, ana.sep)

	if len(tokens) == 0 {
		return tokz, FieldsFuncErr
	}

	if len(ana.normalizers) > 0 {
		tokz = Normalize(tokens, ana.normalizers)
	}

	tokz = tokz.Without(ana.Stopwords())

	return tokz, nil
}

func (ana *Analyzer) WithSep(sep sep.Func) *Analyzer {
	ana.sep = sep
	return ana
}

func (ana *Analyzer) Keywords() *Analyzer {
	ana.sep = nil
	return ana
}

func (ana *Analyzer) WithoutStopWords() bool {
	return len(ana.stopWords) > 0
}

func (ana *Analyzer) SetStopWords(words []string) *Analyzer {
	ana.stopWords = words
	return ana
}

func (ana *Analyzer) Stopwords() AllTokens {
	if !ana.WithoutStopWords() {
		return AllTokens{}
	}
	toks := Normalize(ana.stopWords, ana.normalizers)
	return toks
}

func (ana *Analyzer) IsStopWord(token string) bool {
	return slices.Contains(ana.stopWords, token)
}
