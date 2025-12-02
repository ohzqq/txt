package txt

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/crazy3lf/colorconv"
	"github.com/go-fonts/dejavu/dejavusans"
	"github.com/go-fonts/dejavu/dejavusansmono"
	"github.com/go-fonts/dejavu/dejavuserif"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

var (
	GoMono         = gomono.TTF
	GoRegular      = goregular.TTF
	DejaVuSans     = dejaVuSans()
	DejaVuSansMono = dejaVuSansMono()
	DejaVueSerif   = dejaVuSerif()
)

type Font struct {
	fgColor      string
	bgColor      string
	Width        int
	Height       int
	fontSize     int
	DPI          int
	WrapLines    bool
	paginate     bool
	LineHeight   int
	linesPerPage int
	MaxLines     int
	Wrapper      *Wrapper
	renderer     *etxt.Renderer
}

func NewFont(fnt *sfnt.Font, fs int, opts ...FontOpt) (*Font, error) {
	f := &Font{
		fgColor:    "#000000",
		DPI:        72,
		fontSize:   fs,
		LineHeight: 1,
		MaxLines:   1,
		renderer:   etxt.NewRenderer(),
	}
	f.SetFont(fnt).SetFontSize(fs)
	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}

func ParseFont(src []byte, opts *opentype.FaceOptions) (*Font, error) {
	return &Font{
		fgColor: "#000000",
	}, nil
}

func (f *Font) ParseFont(src []byte, fs int) error {
	fnt, err := NewSFNT(src)
	if err != nil {
		return err
	}
	f.renderer.SetFont(fnt)
	f.SetFontSize(fs)
	return nil
}

func (f *Font) WrapText(str string) []string {
	lines, h := WrapText(str, f.Face(), f.Width)
	f.LineHeight = h
	return lines
}

func (f *Font) SetColors(fg, bg string) *Font {
	return f.SetFgColor(fg).SetBgColor(bg)
}

func (f *Font) SetBgColor(bg string) *Font {
	f.bgColor = bg
	return f
}

func (f *Font) SetFgColor(fg string) *Font {
	f.fgColor = fg
	f.renderer.SetColor(f.getFgColor())
	return f
}

func (f *Font) LinesPerPage() int {
	if f.Height > 0 {
		return f.Height / f.LineHeight
	}
	return f.MaxLines
}

func (f *Font) WriteString(str string) {
	lines := []string{str}
	if f.WrapLines {
		lines = f.WrapText(str)
	}
	pages := NewPaginator(lines, f.LinesPerPage())
	for _, page := range pages.AllPages() {
		txt := strings.TrimSpace(strings.Join(page, "\n"))
		fmt.Printf("%#v\n", txt)
	}
}

func (f *Font) NewBg(w, h int) image.Image {
	return NewImg(f.Width, f.Height, f.getFgColor())
}

func (f *Font) bgSize(text string) (int, int) {
	var rect fract.Rect
	if f.WrapLines {
		rect = f.renderer.MeasureWithWrap(text, f.Width)
	} else {
		rect = f.renderer.Measure(text)
	}
	margin := f.Margin()
	return rect.Max.X.ToInt() + margin*2, rect.Max.Y.ToInt() + margin*2
}

func (f *Font) getFgColor() color.Color {
	fg, err := colorconv.HexToColor(f.fgColor)
	if err == nil {
		return color.Black
	}
	return fg
}

func (f *Font) Margin() int {
	return f.LineHeight - int(f.fontSize)
}

func (f *Font) SetFont(face *sfnt.Font) *Font {
	f.renderer.SetFont(face)
	return f
}

func (f *Font) Face() font.Face {
	fc, _ := opentype.NewFace(f.renderer.GetFont(), f.opentypeOpts())
	return fc
}

func (f *Font) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(f.fontSize),
		DPI:     float64(f.DPI),
		Hinting: font.HintingNone,
	}
}

func (f *Font) SetFontSize(fs int) *Font {
	f.fontSize = fs
	f.renderer.SetSize(float64(fs))
	return f
}

func (f *Font) SetWidth(w int) *Font {
	f.Width = w
	return f
}

func (f *Font) SetSize(w, h int) *Font {
	f.Width = w
	return f.SetHeight(h)
}

func (f *Font) SetHeight(m int) *Font {
	f.Height = m
	return f
}

func dejaVuSans() *sfnt.Font {
	fnt, _ := NewSFNT(dejavusans.TTF)
	return fnt
}

func dejaVuSerif() *sfnt.Font {
	fnt, _ := NewSFNT(dejavuserif.TTF)
	return fnt
}

func dejaVuSansMono() *sfnt.Font {
	fnt, _ := NewSFNT(dejavusansmono.TTF)
	return fnt
}

type FontOpt func(f *Font)

func WithMaxLines(ml int) FontOpt {
	return func(wr *Font) {
		wr.WrapLines = true
		wr.MaxLines = ml
	}
}

func WithLineWrap() FontOpt {
	return func(wr *Font) {
		wr.WrapLines = true
	}
}

func WithSimpleLineWrap(w int) FontOpt {
	return func(wr *Font) {
		wr.WrapLines = true
		wr.Width = w
	}
}

func WithPagination() FontOpt {
	return func(wr *Font) {
		wr.paginate = true
	}
}

func WithSize(w, h int) FontOpt {
	return func(wr *Font) {
		wr.Width = w
		wr.Height = h
		wr.WrapLines = true
	}
}

func WithWidth(w int) FontOpt {
	return func(wr *Font) {
		wr.Width = w
		wr.WrapLines = true
	}
}

func WithHeight(h int) FontOpt {
	return func(wr *Font) {
		wr.Height = h
		wr.WrapLines = true
	}
}

func NewSFNT(src []byte) (*sfnt.Font, error) {
	return sfnt.Parse(src)
}

//The MIT License (MIT)

//Copyright (c) 2012 Grigory Dryapak

//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

func NewImg(width, height int, fillColor color.Color) *image.NRGBA {
	if width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
	if (c == color.NRGBA{0, 0, 0, 0}) {
		return image.NewNRGBA(image.Rect(0, 0, width, height))
	}

	return &image.NRGBA{
		Pix:    bytes.Repeat([]byte{c.R, c.G, c.B, c.A}, width*height),
		Stride: 4 * width,
		Rect:   image.Rect(0, 0, width, height),
	}
}
