package txt

import (
	"errors"
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

var (
	FieldsFuncErr = errors.New("strings.FieldsFunc returned an empty slice or the string was empty")
	EmptyStrErr   = errors.New("empty string")
)

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

func (ana *Analyzer) Tokenize(text string) (*Field, error) {
	var (
		f   = &Field{}
		err error
	)
	f.Tokens, err = Tokenize(ana, text)
	if err != nil {
		return f, err
	}

	return f, nil
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
