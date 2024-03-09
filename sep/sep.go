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

const (
	Whitespace = iota
	Space
	Newline
	Comma
	Tab
)

func Split(str string, seps ...Sep) []string {
	if len(seps) == 0 {
		seps = append(seps, Space)
	}
	fields := splitField(str, seps[0])
	for _, sep := range seps[0:] {
		for _, s := range fields {
			fields = append(fields, Split(s, sep)...)
		}
	}
	return fields
}

func splitField(sep Sep, fields ...string) []string {
	var split []string
	switch len(fields) {
	case 0:
		return split
	case 1:
		split = append(split, strings.FieldsFunc(fields[0], sep))
	default:
		split = append(splits, splitField)
	}
}

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
