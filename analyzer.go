package txt

import "strings"

type Analyzer struct {
	StopWords    []string
	Stem         bool
	AlphaNumOnly bool
	ToLower      bool
}

type Opt func(*Analyzer)

func NewAnalyzer(opts ...Opt) *Analyzer {
	ana := &Analyzer{
		AlphaNumOnly: true,
		ToLower:      true,
	}
	for _, opt := range opts {
		opt(ana)
	}
	return ana
}

func (ana *Analyzer) RmStopWords() bool {
	return len(ana.StopWords) > 0
}

func NormalizeText(token string) string {
	fields := toLower(strings.Split(token, " "))
	for t, term := range fields {
		if len(term) == 1 {
			fields[t] = term
		} else {
			fields[t] = stripAlphaNum(term)
		}
	}
	return strings.Join(fields, " ")
}

func KeepCase(ana *Analyzer) {
	ana.ToLower = false
}

func toLower(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func KeepPunct(ana *Analyzer) {
	ana.AlphaNumOnly = false
}

func stripAlphaNum(token string) string {
	var s []byte
	for _, b := range []byte(token) {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s = append(s, b)
		}
	}
	return string(s)
}

func Stem(ana *Analyzer) {
	ana.Stem = true
}

func stemWords(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = english.Stem(token, false)
	}
	return r
}

func WithStopWords(words []string) Opt {
	return func(ana *Analyzer) {
		ana.StopWords = words
	}
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
