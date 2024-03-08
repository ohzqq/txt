package txt

import (
	"errors"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
)

type Analyzer struct {
	StopWords    []string
	Stem         bool
	AlphaNumOnly bool
	ToLower      bool
	fieldsFunc   func(r rune) bool
	normalizers  []Normalizer
}

type Normalizer func(string) string

type Opt func(*Analyzer)

func NewAnalyzer(normalizers ...Normalizer) *Analyzer {
	ana := &Analyzer{
		fieldsFunc:  func(r rune) bool { return unicode.IsSpace(r) },
		normalizers: normalizers,
	}
	return ana
}

func (ana *Analyzer) Tokenize(text string) ([]*Token, error) {
	var (
		toks   []*Token
		tokens = strings.FieldsFunc(text, ana.fieldsFunc)
	)

	if len(tokens) == 0 {
		return toks, errors.New("strings.FieldsFunc returned an empty slice or the string was empty")
	}

	for _, label := range tokens {
		tok := label
		for _, norm := range ana.normalizers {
			tok = norm(tok)
		}

		if ana.RmStopWords() {
			if ana.IsStopWord(tok) {
				continue
			}
		}

		if tok != "" {
			toks = append(toks, NewToken(label, tok))
		}
	}

	return toks, nil
}

func (ana *Analyzer) RmStopWords() bool {
	return len(ana.StopWords) > 0
}

func (ana *Analyzer) SetStopWords(words []string) {
	ana.StopWords = words
}

func (ana *Analyzer) IsStopWord(token string) bool {
	return slices.Contains(ana.StopWords, token)
}

func Normalize(ana *Analyzer) {
	ana.AlphaNumOnly = true
	ana.ToLower = true
}

func normalizeText(token string) string {
	fields := toLower(strings.Split(token, " "))
	for t, term := range fields {
		if len(term) == 1 {
			fields[t] = term
		} else {
			fields[t] = RmPunct(term)
		}
	}
	return strings.Join(fields, " ")
}

func ToLower(ana *Analyzer) {
	ana.ToLower = true
}

func toLower(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func NoPunct(ana *Analyzer) {
	ana.AlphaNumOnly = true
}

func KeepPunct(ana *Analyzer) {
	ana.AlphaNumOnly = false
}

func RmPunct(token string) string {
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

func StemWord(ana *Analyzer) {
	ana.Stem = true
}

func stemWords(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = english.Stem(token, false)
	}
	return r
}

func Stem(token string) string {
	return english.Stem(token, false)
}

func WithStopWords(words []string) Opt {
	return func(ana *Analyzer) {
		ana.StopWords = words
	}
}

func (ana *Analyzer) RmStopWord(word string) string {
	if !lo.Contains(ana.StopWords, word) {
		return word
	}
	return ""
}

func rmStopWords(tokens []string) []string {
	var words []string
	for _, token := range tokens {
		if !lo.Contains(stopWords, token) {
			words = append(words, token)
		}
	}
	return words
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
