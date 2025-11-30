package txt

import (
	"testing"
)

func TestVariableTextWrapperPaginator(t *testing.T) {
	wr := NewTextWrapper()
	wr.
		SetFontSize(12).
		SetWidth(250).
		WithGoMono()
	box := NewPaginator(wr).SetHeight(90)
	pages := box.Paginate(tstStr)
	pages.SetLoop(false)
}

func TestSimplePaginator(t *testing.T) {
	wr := NewSimpleTextWrapper(250)
	pagy := NewPaginator(wr).SetMaxLines(3)
	pages := pagy.Paginate(tstStr)
	println(pages.CurrentPage())
}

func TestTextWrap(t *testing.T) {
	//lines := NewFrame(tstStr, 26, 250)
	box := NewTextWrapper()
	box.
		SetFontSize(2).
		SetWidth(250).
		WithGoMono()
	//err := box.WithTTF(GoMono)
	//if err != nil {
	//t.Error(err)
	//}
	box.WrapText(tstStr)
	//lines := box.WrapText(tstStr)
	//fmt.Printf("%#v\n", pages)
}

const tstStr = "My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around."
