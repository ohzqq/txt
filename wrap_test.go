package txt

import (
	"testing"
)

func TestTextWrap(t *testing.T) {
	//lines := NewFrame(tstStr, 26, 250)
	box := NewWrapper()
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
