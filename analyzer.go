package txt

import (
	"slices"
	"strings"

	"github.com/ohzqq/txt/sep"
)

type Analyzer struct {
	StopWords   []string
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

func Normalize(opts ...Option) *Analyzer {
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
	ana.StopWords = stopWords
}

func ToLower(ana *Analyzer) {
	ana.normalizers = append(ana.normalizers, strings.ToLower)
}

func WithoutPunct(ana *Analyzer) {
	ana.normalizers = append(ana.normalizers, StripPunct)
}

func WithStemmer(ana *Analyzer) {
	ana.normalizers = append(ana.normalizers, Stem)
}

func WithNormalizers(n ...Normalizer) Option {
	return func(ana *Analyzer) {
		ana.normalizers = append(ana.normalizers, Stem)
	}
}

func (ana *Analyzer) WithSep(sep sep.Func) *Analyzer {
	ana.sep = sep
	return ana
}

func (ana *Analyzer) Keywords() *Analyzer {
	ana.sep = nil
	return ana
}

func (ana *Analyzer) RmStopWords() bool {
	return len(ana.StopWords) > 0
}

func (ana *Analyzer) SetStopWords(words []string) *Analyzer {
	ana.StopWords = words
	return ana
}

func (ana *Analyzer) IsStopWord(token string) bool {
	return slices.Contains(ana.StopWords, token)
}
