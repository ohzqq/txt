package txt

import (
	"slices"
	"strings"

	"github.com/kljensen/snowball/english"
	"github.com/ohzqq/txt/sep"
)

type Analyzer struct {
	StopWords   []string
	sep         sep.Func
	normalizers []Normalizer
}

type Normalizer func(string) string

type Option func(*Analyzer)

func New(normies ...Normalizer) *Analyzer {
	return &Analyzer{
		normalizers: normies,
		sep:         sep.Whitespace,
	}
}

func Keywords() *Analyzer {
	return New(strings.ToLower).Keywords()
}

func Normalize(opts ...Option) *Analyzer {
	ana := New(
		strings.ToLower,
		AlphaNum,
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

func WithStemmer(ana *Analyzer) {
	ana.normalizers = append(ana.normalizers, Stem)
}

func (ana *Analyzer) Tokenize(text string) (Tokens, error) {
	return Tokenize(text, ana)
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

func AlphaNum(token string) string {
	var s []byte
	for _, b := range []byte(token) {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			s = append(s, b)
		}
	}
	return string(s)
}

func Stem(token string) string {
	return english.Stem(token, false)
}

func DefaultStopWords() []string {
	return stopWords
}

var stopWords = []string{
	"i",
	"vol",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"of",
	"at",
	"by",
	"for",
	"with",
	"into",
	"to",
	"from",
	"then",
	"when",
	"where",
	"why",
	"how",
	"no",
	"not",
	"than",
	"too",
}
