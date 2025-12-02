//go:build gtxt

package txt

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"
)

var (
	fontSize    = 18
	testBoxSize = WithSize(250, 100)
)

func TestTextWrap(t *testing.T) {
	t.Skip()
	//lines := NewFrame(tstStr, 26, 250)
	//lines, totalLines := WrapFont(tstStr, WithGoMono(22), WithSize(250, 100))
	ff, err := NewFont(DejaVuSans, 22, WithSimpleLineWrap(5))
	if err != nil {
		t.Fatal(err)
	}
	lines := ff.WrapText(tstStr)
	linesPerPage := ff.LinesPerPage()
	//lines, totalLines := WrapTextbox(tstStr, 250, 100)
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
	fmt.Printf("%#v\n", linesPerPage)
}

func TestPagination(t *testing.T) {
	ff, err := NewFont(DejaVuSans, fontSize, testBoxSize)
	if err != nil {
		t.Fatal(err)
	}
	lines := ff.WrapText(tstStr)
	want := ff.LinesPerPage()
	pages := NewPaginator(lines, want)
	for _, page := range pages.AllPages() {
		if got := len(page); got > want {
			t.Errorf("got %v lines, wanted %v\n", got, want)
		}
	}
}

func TestMaxLinesPerPage(t *testing.T) {
	ff, err := NewFont(DejaVuSans, fontSize, WithMaxLines(3))
	if err != nil {
		t.Fatal(err)
	}
	ff.WrapText(tstStr)
	want := 3
	if got := ff.LinesPerPage(); got != want {
		t.Errorf("got %v lines, wanted %v\n", got, want)
	}
}

func TestHeightLinesPerPage(t *testing.T) {
	ff, err := NewFont(DejaVuSans, fontSize, testBoxSize)
	if err != nil {
		t.Fatal(err)
	}
	ff.WrapText(tstStr)
	want := 4
	if got := ff.LinesPerPage(); got != want {
		t.Errorf("got %v lines, wanted %v\n", got, want)
	}
}

func TestEtxt(t *testing.T) {
	ff, err := NewFont(DejaVuSans, fontSize, testBoxSize)
	if err != nil {
		t.Fatal(err)
	}
	lines := ff.WrapText(tstStr)
	pages := NewPaginator(lines, ff.LinesPerPage())
	txt := pages.JoinCurrentWithSpace()
	w, h := ff.bgSize(txt)
	dst := NewImg(w, h, color.White)
	lh := -ff.Margin()
	for _, line := range pages.Current() {
		lh = lh + ff.LineHeight
		ff.renderer.Draw(dst, line, ff.Margin(), lh)
	}

	out := `testdata/test.png`
	testImg, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	err = png.Encode(testImg, dst)
	if err != nil {
		t.Fatal(err)
	}
}

const tstStr = "My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around."
