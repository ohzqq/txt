package txt

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/kljensen/snowball/english"
	"github.com/spf13/cast"
)

type Keyword struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Children *Field
	bits     *roaring.Bitmap
}

func NewKeyword(label string) *Keyword {
	return &Keyword{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (kw *Keyword) Bitmap() *roaring.Bitmap {
	return kw.bits
}

func (kw *Keyword) SetValue(txt string) *Keyword {
	kw.Value = txt
	return kw
}

func (kw *Keyword) Items() []int {
	i := kw.bits.ToArray()
	return cast.ToIntSlice(i)
}

func (kw *Keyword) Count() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Keyword) Len() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Keyword) Contains(id int) bool {
	return kw.bits.ContainsInt(id)
}

func (kw *Keyword) Add(ids ...int) {
	for _, id := range ids {
		if !kw.Contains(id) {
			kw.bits.AddInt(id)
		}
	}
}

func KeywordTokenizer(val any) []*Keyword {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Keyword, len(tokens))
	for i, token := range tokens {
		items[i] = NewKeyword(token)
		items[i].Value = normalizeText(token)
	}
	return items
}

func NormalizeText(token string) string {
	fields := lowerCase(strings.Split(token, " "))
	for t, term := range fields {
		if len(term) == 1 {
			fields[t] = term
		} else {
			fields[t] = stripNonAlphaNumeric(term)
		}
	}
	return strings.Join(fields, " ")
}

func LowerCase(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func StripNonAlphaNumeric(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}

func Stem(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = english.Stem(token, false)
	}
	return r
}

func StopWords(tokens []string) []string {
	var words []string
	for _, token := range tokens {
		if !lo.Contains(stopWords, token) {
			words = append(words, token)
		}
	}
	return words
}
