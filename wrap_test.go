//go:build gtxt

package txt

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
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

func TestEtxt(t *testing.T) {
	fs := 40
	ff := NewFont(WithGoRegular(fs), WithLineWrap(), WithMaxLines(3))
	ff.SetWidth(250)
	lines := ff.WrapText(tstStr)
	pages := NewPaginator(lines, ff.LinesPerPage())
	fmt.Printf("%#v\n", pages.Current())
	//ff := NewFont(WithGoMono(22))
	//ff.calculateBounds(tstStr)
	//face, err := sfnt.Parse(goRegular)
	//if err != nil {
	//  t.Fatal(err)
	//}
	//ctx := etxt.NewRenderer()
	//ctx.SetFont(face)
	//ctx.SetColor(color.Black)
	//ctx.SetSize(float64(fs))
	//utils := ctx.Utils()
	//lh := utils.GetLineHeight()
	txt := strings.TrimSpace(strings.Join(pages.Current(), " "))
	w, h := ff.bgSize(txt)
	//fmt.Printf("lineheight: %d\ntotal height: %d\ntotal lines: %d\n", int(lh), size.Max.Y.ToInt(), size.Max.Y.ToInt()/ff.Wrapper.LineHeight)
	dst := NewImg(w, h, color.White)
	lh := -ff.Margin()
	for _, line := range pages.Current() {
		println(line)
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
