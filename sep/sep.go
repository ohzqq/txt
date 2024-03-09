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
	Whitespace Sep = iota
	Space
	Newline
	Comma
	Tab
)

func Split(str string, seps ...Sep) []string {
	if len(seps) == 0 {
		return []string{str}
	}
	return strings.FieldsFunc(str, separators[seps[0]])
}

func ShouldSplit(seps []Sep) func(r rune) bool {
	return func(r rune) bool {
		for _, sep := range seps {
			if sh := separators[sep](r); sh {
				return sh
			}
		}
		return false
	}
}
