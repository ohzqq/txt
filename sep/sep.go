package sep

import (
	"strings"
	"unicode"
)

type Func func(r rune) bool

func Split(str string, seps ...Func) []string {
	if len(seps) == 0 {
		return []string{str}
	}
	return strings.FieldsFunc(str, seps[0])
}

func Whitespace(str string) []string {
	return strings.FieldsFunc(str, OnWhitespace)
}

func Punct(str string) []string {
	return strings.FieldsFunc(str, OnPunct)
}

func Comma(str string) []string {
	return strings.FieldsFunc(str, OnComma)
}

func Newline(str string) []string {
	return strings.FieldsFunc(str, OnNewline)
}

func Tab(str string) []string {
	return strings.FieldsFunc(str, OnTab)
}

func OnPunct(r rune) bool      { return unicode.IsPunct(r) }
func OnComma(r rune) bool      { return r == ',' }
func OnWhitespace(r rune) bool { return unicode.IsSpace(r) }
func OnSpace(r rune) bool      { return r == ' ' }
func OnTab(r rune) bool        { return r == '\t' }
func OnNewline(r rune) bool    { return r == '\r' || r == '\n' }
