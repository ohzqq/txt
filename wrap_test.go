package txt

import (
	"fmt"
	"testing"
)

func TestTextWrap(t *testing.T) {
	//lines := NewFrame(tstStr, 26, 250)
	//lines, totalLines := WrapFont(tstStr, WithGoMono(22), WithSize(250, 100))
	lines, totalLines := WrapTextbox(tstStr, 250, 100, WithGoMono(22))
	//lines, totalLines := SimpleWrap(tstStr, 25, 4)
	//box := NewWrapper()
	//box.
	//SetFontSize(30).
	//SetWidth(250).
	//WithGoMono()
	//err := box.WithTTF(GoMono)
	//if err != nil {
	//t.Error(err)
	//}
	//lines := box.WrapText(tstStr)
	fmt.Printf("total %v: %#v\n", len(lines), lines)
	fmt.Printf("%#v\n", totalLines)
}

func TestPagination(t *testing.T) {
	box := NewWrapper(WithGoMono(22))
	box.
		SetFontSize(30).
		SetWidth(250)
	pages := box.Paginate(tstStr)
	fmt.Printf("%#v\n", pages.AllPages())
}

const tstStr = "My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around."
