package sep

import (
	"strings"
	"unicode"
)

var (
	whitespace = func(r rune) bool { return unicode.IsSpace(r) }
	comma      = func(r rune) bool { return r == ',' }
	space      = func(r rune) bool { return r == ' ' }
	tab        = func(r rune) bool { return r == '\t' }
	newline    = func(r rune) bool { return r == '\r' || r == '\n' }
	separators = map[Sep]func(r rune) bool{
		Whitespace: whitespace,
		Space:      space,
		Comma:      comma,
		Tab:        tab,
		Newline:    newline,
	}
)

type Sep int

type Splitter func(string) []string

const (
	Whitespace Sep = iota
	Space
	Newline
	Comma
	Tab
)

func Split(str string, seps ...Sep) []string {
	if len(seps) == 0 {
		return strings.Split(str, " ")
	}
	//fields = append(fields, strings.FieldsFunc(str, separators[seps[0]])...)
	//c := 1
	//if len(seps) > c {
	//for i := c; i < len(seps); i++ {
	var fields []string
	for _, sep := range seps {
		for _, s := range fields {
			//println(s)
			//fields = append(fields, strings.FieldsFunc(s, separators[seps[i]])...)
			fields = append(fields, splitStrings(sep, s)...)
		}
	}
	return fields
}

func splitStrings(sep Sep, strs ...string) []string {
	var s []string
	for _, str := range strs {
		s = append(s, strings.FieldsFunc(str, separators[sep])...)
	}
	return s
}

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
