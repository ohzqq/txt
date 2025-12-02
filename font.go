//go:build gtxt

package txt

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"strings"

	"github.com/crazy3lf/colorconv"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var (
	goMono             = gomono.TTF
	goRegular          = goregular.TTF
	inconsolataRegular = inconsolata.Regular8x16
)

type Font struct {
	fgColor      string
	bgColor      string
	bg           draw.Image
	fg           image.Image
	Width        int
	Height       int
	FontSize     int
	DPI          int
	WrapLines    bool
	paginate     bool
	LineHeight   int
	linesPerPage int
	MaxLines     int
	Wrapper      *Wrapper
	font         *sfnt.Font
	renderer     *etxt.Renderer
	//font *Face
}

func NewFont(opts ...FontOpt) *Font {
	f := &Font{
		fgColor:    "#000000",
		FontSize:   16,
		DPI:        72,
		LineHeight: 1,
		MaxLines:   1,
		Wrapper:    NewWrapper(),
		renderer:   etxt.NewRenderer(),
	}
	for _, opt := range opts {
		opt(f)
	}
	if f.font == nil {
		f.ParseFont(goRegular, f.FontSize)
	}
	f.renderer.SetFont(f.font)
	f.renderer.SetColor(color.Black)
	f.Wrapper.Height = f.Height
	return f
}

func ParseFont(src []byte, opts *opentype.FaceOptions) (*Font, error) {
	return &Font{
		fgColor: "#000000",
	}, nil
}

func (f *Font) WrapText(str string) []string {
	return f.Wrapper.WrapText(str, f)
}

func (f *Font) SetColors(fg, bg string) *Font {
	return f.SetFgColor(fg).SetBgColor(bg)
}

func (f *Font) SetBgColor(bg string) *Font {
	f.bgColor = bg
	return f
}

func (f *Font) SetFgColor(fg string) *Font {
	fgc := f.getFgColor()
	f.renderer.SetColor(fgc)
	f.fgColor = fg
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
	pages := NewPaginator(lines, f.Wrapper.LinesPerPage())
	for _, page := range pages.AllPages() {
		txt := strings.TrimSpace(strings.Join(page, "\n"))
		fmt.Printf("%#v\n", txt)
	}
}

func (f *Font) Face() font.Face {
	fc, _ := opentype.NewFace(f.font, f.opentypeOpts())
	return fc
}

func (f *Font) WriteStringTo(wr io.Writer, str string) error {
	lines := []string{str}
	if f.WrapLines {
		lines = f.WrapText(str)
	}
	pages := NewPaginator(lines, f.Wrapper.LinesPerPage())
	for _, page := range pages.AllPages() {
		txt := strings.TrimSpace(strings.Join(page, "\n"))
		_, err := wr.Write([]byte(txt))
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Font) DrawString(str string) *font.Drawer {
	drawer := f.GetDrawer()
	drawer.DrawString(str)
	return drawer
}

func (f *Font) GetDrawer() *font.Drawer {
	f.NewBg(f.Width, f.Height)
	return &font.Drawer{
		Dst: f.bg,
		Src: image.NewUniform(f.getFgColor()),
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fixed.I(0),
		},
	}
}

func (f *Font) NewBg(w, h int) *Font {
	//var bc color.Color
	//var err error
	//bc = color.Transparent
	//if f.bgColor != "" {
	//  bc, err = colorconv.HexToColor(f.bgColor)
	//  if err != nil {
	//  }
	//}
	img := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	//println(f.bgColor)
	return f.SetBg(img)
}

//func (f *Font) Face() font.Face {
//}

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

func (f *Font) SetBg(img draw.Image) *Font {
	f.bg = img
	return f
}

func (f *Font) getFgColor() color.Color {
	fg, err := colorconv.HexToColor(f.fgColor)
	if err == nil {
		return color.Black
	}
	return fg
}

func (f *Font) Margin() int {
	return f.LineHeight - f.FontSize
}

func (f *Font) SetFont(face *sfnt.Font) *Font {
	f.font = face
	return f
}

func (f *Font) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(f.FontSize),
		DPI:     float64(f.DPI),
		Hinting: font.HintingNone,
	}
}

func (f *Font) ParseFont(src []byte, fs int) error {
	fnt, err := NewSFNT(src)
	if err != nil {
		return err
	}
	f.font = fnt
	f.SetFontSize(fs)
	return nil
}

func (f *Font) ParseTTF(src []byte, fs int) error {
	fnt, err := NewSFNT(goRegular)
	if err != nil {
		return err
	}
	f.SetFontSize(fs)
	f.SetFont(fnt)
	return nil
}

func (f *Font) SetFontSize(fs int) *Font {
	f.FontSize = fs
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
	f.Wrapper.Height = f.Height
	return f
}

type FontOpt func(f *Font)

func WithFont(face *sfnt.Font) FontOpt {
	return func(f *Font) {
		f.SetFont(face)
	}
}

func WithMaxLines(ml int) FontOpt {
	return func(wr *Font) {
		wr.MaxLines = ml
		wr.Wrapper.MaxLines = ml
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
		wr.Wrapper.Simple = true
		wr.Width = w
	}
}

func WithPagination() FontOpt {
	return func(wr *Font) {
		wr.paginate = true
	}
}

func WithGoMono(fs int) FontOpt {
	return func(wr *Font) {
		wr.ParseFont(goMono, fs)
	}
}

func WithGoRegular(fs int) FontOpt {
	return func(wr *Font) {
		wr.ParseFont(goRegular, fs)
	}
}
func WithSize(w, h int) FontOpt {
	return func(wr *Font) {
		wr.Width = w
		wr.Height = h
	}
}

func WithFontSize(fs int) FontOpt {
	return func(wr *Font) {
		wr.FontSize = fs
	}
}

//func GoRegular(opts *opentype.FaceOptions) *Font {
//  f, _ := NewTTF(goRegular, opts)
//  return &Font{
//    Face: f,
//  }
//}

//func GoMono(opts *opentype.FaceOptions) *Font {
//  f, _ := NewTTF(goMono, opts)
//  return &Font{
//    Face: f,
//  }
//}

func Inconsolata() *Font {
	return &Font{}
}

func NewSFNT(src []byte) (*sfnt.Font, error) {
	return sfnt.Parse(src)
}

func NewTTF(src []byte, opts *opentype.FaceOptions) (font.Face, error) {
	inc, err := opentype.Parse(src)
	if err != nil {
		return &opentype.Face{}, err
	}
	return opentype.NewFace(inc, opts)
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
