package txt

import (
	"fmt"
	"image"
	"testing"

	wordwrap "github.com/arran4/golang-wordwrap"
)

func TestWordWrap(t *testing.T) {
	f, err := NewTTF("testdata/Inconsolata-Regular.ttf", 26)
	if err != nil {
		t.Fatal(err)
	}
	rect := image.Rect(0, 0, 250, 600)
	_, lines, _, err := wordwrap.SimpleWrapTextToRect(tstStr, rect, f)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range lines {
		println(line.TextValue())
	}
}

func TestWrapToSlice(t *testing.T) {
	lines := WrapToSlice(tstStr, 40)
	if len(lines) != 6 {
		fmt.Printf("%#v\n", len(lines))
	}
}

func TestWrapToSliceChunk(t *testing.T) {
	lines := WrapAndChunk(tstStr, 40, 3)
	if len(lines) != 2 {
		fmt.Printf("%#v\n", len(lines))
	}
}

func TestTextbox(t *testing.T) {
	//lines := NewFrame(tstStr, 26, 250)
	box := NewTextWrapper()
	box.SetWidth(250).SetMaxLines(0)
	err := box.WithTTF(GoMono)
	if err != nil {
		t.Error(err)
	}
	pages := box.Paginate(tstStr)
	//lines := box.WrapText(tstStr)
	fmt.Printf("%#v\n", pages)
}

const tstStr = "My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around."
