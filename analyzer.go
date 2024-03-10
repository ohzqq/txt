package txt

import (
	"slices"
	"strings"

	"github.com/ohzqq/txt/sep"
)

type Analyzer struct {
	stopWords   []string
	sep         sep.Func
	normalizers []Normalizer
}

type Option func(*Analyzer)

func New(opts ...Option) *Analyzer {
	ana := &Analyzer{
		sep: sep.Whitespace,
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

func WithStopWords(words []string) Option {
	return func(ana *Analyzer) {
		ana.SetStopWords(words)
	}
}

func WithDefaultStopWords(ana *Analyzer) {
	ana.SetStopWords(stopWords)
}

func ToLower(ana *Analyzer) {
	ana.AddNormalizer(strings.ToLower)
}

func WithoutPunct(ana *Analyzer) {
	ana.AddNormalizer(StripPunct)
}

func WithStemmer(ana *Analyzer) {
	ana.AddNormalizer(Stem)
}

func WithNormalizers(n ...Normalizer) Option {
	return func(ana *Analyzer) {
		ana.AddNormalizer(n...)
	}
}

func (ana *Analyzer) WithNormalizer(normies ...Normalizer) *Analyzer {
	ana.normalizers = normies
	return ana
}

func (ana *Analyzer) AddNormalizer(normies ...Normalizer) *Analyzer {
	ana.normalizers = append(ana.normalizers, normies...)
	return ana
}

func (ana *Analyzer) Tokenize(text string) (Tokens, error) {
	var (
		toks   Tokens
		tokens []string
	)

	if text == "" {
		return toks, EmptyStrErr
	}

	tokens = Split(text, ana.sep)

	if len(tokens) == 0 {
		return toks, FieldsFuncErr
	}

	if len(ana.normalizers) > 0 {
		toks = Normalize(tokens, ana.normalizers)
	}

	toks = toks.Without(ana.Stopwords())

	return toks, nil
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

func (ana *Analyzer) Stopwords() Tokens {
	if !ana.WithoutStopWords() {
		return Tokens{}
	}
	toks := Normalize(ana.stopWords, ana.normalizers)
	return toks
}

func (ana *Analyzer) IsStopWord(token string) bool {
	return slices.Contains(ana.stopWords, token)
}
