package txt

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
)

func TestTextWrap(t *testing.T) {
	//lines := NewFrame(tstStr, 26, 250)
	//lines, totalLines := WrapFont(tstStr, WithGoMono(22), WithSize(250, 100))
	ff := NewFont(WithGoMono(22), WithSimpleLineWrap(5))
	lines := ff.WrapText(tstStr)
	linesPerPage := ff.Wrapper.LinesPerPage()
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
	ff := NewFont(WithGoMono(22), WithLineWrap(), WithSize(250, 100))
	lines := ff.WrapText(tstStr)
	pages := NewPaginator(lines, ff.Wrapper.LinesPerPage())
	fmt.Printf("%#v\n", pages.AllPages())
}

func TestWriteFont(t *testing.T) {
	ff := NewFont(WithGoMono(22))
	err := ff.WriteStringTo(os.Stdout, tstStr)
	if err != nil {
		t.Fatal(err)
	}
	//pages := NewPaginator(lines, ff.Wrapper.LinesPerPage())
	//fmt.Printf("%#v\n", pages.AllPages())
}

func TestEncodeImg(t *testing.T) {
	ff := NewFont(WithGoMono(22))
	ff.calculateBounds(tstStr)
	//ff.bgColor = "#ffffff"
	//dr := ff.DrawString(tstStr)
	out := `testdata/test.png`
	//format := imgconv.FormatOption{Format: imgconv.PNG}
	//err := imgconv.Save(out, dr.Dst, format)
	ttf, err := os.Open("testdata/Inconsolata-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	defer ttf.Close()
	face, err := font.ParseTTF(ttf)
	if err != nil {
		t.Fatal(err)
	}

	img := NewImg(ff.Width, ff.Height, color.White)
	renderer := &render.Renderer{
		FontSize: float32(ff.FontSize),
		Color:    color.Black,
	}
	renderer.DrawString(tstStr, img, face)
	testImg, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	err = png.Encode(testImg, img)
	if err != nil {
		t.Fatal(err)
	}
}

const tstStr = "My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around."
