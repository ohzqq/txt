package sep

import (
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

type Seps []Sep

const (
	Whitespace Sep = iota
	Space
	Newline
	Comma
	Tab
)

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

func (s Sep) SplitOn(r rune) bool {
	return separators[s](r)
	//switch s {
	//case Whitespace:
	//  return unicode.IsSpace(r)
	//case Newline:
	//  return r == '\r' || r == '\n'
	//case Comma:
	//  return r == ','
	//case Tab:
	//  return r == '\t'
	//case Space:
	//  return r == ' '
	//default:
	//  return false
	//}

}

//func (s Sep) Contains(s string, r rune) bool {
//}

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
