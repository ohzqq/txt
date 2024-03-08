package txt

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	bits  *roaring.Bitmap
}

func NewToken(label, val string) *Token {
	return &Token{
		Value: val,
		Label: label,
		bits:  roaring.New(),
	}
}

func (kw *Token) Bitmap() *roaring.Bitmap {
	return kw.bits
}

func (kw *Token) SetValue(txt string) *Token {
	kw.Value = txt
	return kw
}

func (kw *Token) Items() []int {
	i := kw.bits.ToArray()
	return cast.ToIntSlice(i)
}

func (kw *Token) Count() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Len() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Contains(id int) bool {
	return kw.bits.ContainsInt(id)
}

func (kw *Token) Add(ids ...int) {
	for _, id := range ids {
		if !kw.Contains(id) {
			kw.bits.AddInt(id)
		}
	}
}

func KeywordTokenizer(val any) []*Token {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Token, len(tokens))
	for i, token := range tokens {
		val := normalizeText(token)
		items[i] = NewToken(token, val)
	}
	return items
}
