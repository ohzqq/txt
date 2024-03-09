package txt

import (
	"errors"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

var (
	SpaceSep   = func(r rune) bool { return unicode.IsSpace(r) }
	CommaSep   = func(r rune) bool { return r == ',' }
	TabSep     = func(r rune) bool { return r == '\t' }
	NewlineSep = func(r rune) bool { return r == '\r' || r == '\n' }
)

type Analyzer struct {
	StopWords   []string
	fieldsFunc  func(r rune) bool
	normalizers []Normalizer
}

type Normalizer func(string) string

type Sep func(r rune) bool

func NewAnalyzer(normalizers []Normalizer, sep ...Sep) *Analyzer {
	ana := &Analyzer{
		normalizers: normalizers,
	}
	if len(sep) > 0 {
		ana.fieldsFunc = sep[0]
	}
	return ana
}

func Simple() *Analyzer {
	return NewAnalyzer([]Normalizer{})
}

func Keyword(normalizers ...Normalizer) *Analyzer {
	return NewAnalyzer(normalizers)
}

func SplitOnSpaces(normalizers ...Normalizer) *Analyzer {
	ana := NewAnalyzer(normalizers, SpaceSep)
	return ana
}

func Normalize(splitter ...Sep) *Analyzer {
	ana := NewAnalyzer([]Normalizer{
		strings.ToLower,
		AlphaNum,
	}, splitter...)
	return ana
}

func Complex(splitter ...Sep) *Analyzer {
	ana := NewAnalyzer([]Normalizer{
		strings.ToLower,
		AlphaNum,
		Stem,
	}, splitter...)
	ana.StopWords = DefaultStopWords()
	return ana
}

func (ana *Analyzer) Tokenize(text string) (Tokens, error) {
	var (
		toks   Tokens
		tokens []string
	)

	if ana.fieldsFunc == nil {
		tokens = []string{text}
	} else {
		tokens = strings.FieldsFunc(text, ana.fieldsFunc)
	}

	if len(tokens) == 0 {
		return toks, errors.New("strings.FieldsFunc returned an empty slice or the string was empty")
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

func (ana *Analyzer) SetFieldsFunc(fn func(r rune) bool) {
	ana.fieldsFunc = fn
}

func (ana *Analyzer) Keywords() {
	ana.SetFieldsFunc(nil)
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

func NormalizeWord(word string) string {
	word = strings.ToLower(word)
	word = AlphaNum(word)
	return word
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
