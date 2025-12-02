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
	fgColor   string
	bgColor   string
	bg        draw.Image
	fg        image.Image
	Width     int
	Height    int
	FontSize  int
	DPI       int
	WrapLines bool
	paginate  bool
	Wrapper   *Wrapper
	font.Face
	font *sfnt.Font
}

func NewFont(opts ...FontOpt) *Font {
	f := &Font{
		fgColor:  "#000000",
		FontSize: 16,
		DPI:      72,
		Wrapper:  NewWrapper(),
	}
	fnt, _ := NewSFNT(goRegular)
	f.font = fnt
	for _, opt := range opts {
		opt(f)
	}
	if f.Face == nil {
		f.Face = Inconsolata()
	}
	f.Wrapper.Height = f.Height
	return f
}

func ParseFont(src []byte, opts *opentype.FaceOptions) (*Font, error) {
	f, err := NewTTF(src, opts)
	if err != nil {
		return nil, err
	}
	return &Font{
		fgColor: "#000000",
		Face:    f,
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
	f.fgColor = fg
	return f
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

func (f *Font) calculateBounds(str string) {
	if f.Width == 0 {
		rect, _ := font.BoundString(f.Face, str)
		f.Width = rect.Max.X.Round() // 2
	}
	if f.Height == 0 {
		f.Height = f.Face.Metrics().Height.Round()
	}
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
	f.calculateBounds(str)
	drawer := f.GetDrawer()
	drawer.DrawString(str)
	return drawer
}

func (f *Font) GetDrawer() *font.Drawer {
	f.NewBg(f.Width, f.Height)
	return &font.Drawer{
		Dst:  f.bg,
		Src:  image.NewUniform(f.getFgColor()),
		Face: f.Face,
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

func (f *Font) SetFont(face font.Face) *Font {
	f.Face = face
	return f
}

func (f *Font) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(f.FontSize),
		DPI:     float64(f.DPI),
		Hinting: font.HintingNone,
	}
}

func (wr *Font) ParseTTF(src []byte, fs int) error {
	wr.SetFontSize(fs)
	f, err := ParseFont(src, wr.opentypeOpts())
	if err != nil {
		return err
	}
	wr.SetFont(f)
	return nil
}

func (f *Font) SetFontSize(fs int) *Font {
	f.FontSize = fs
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

func WithFont(face font.Face) FontOpt {
	return func(f *Font) {
		f.SetFont(face)
	}
}

func WithMaxLines(ml int) FontOpt {
	return func(wr *Font) {
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
		wr.ParseTTF(goMono, fs)
	}
}

func WithGoRegular(fs int) FontOpt {
	return func(wr *Font) {
		wr.ParseTTF(goRegular, fs)
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

func GoRegular(opts *opentype.FaceOptions) *Font {
	f, _ := NewTTF(goRegular, opts)
	return &Font{
		Face: f,
	}
}

func GoMono(opts *opentype.FaceOptions) *Font {
	f, _ := NewTTF(goMono, opts)
	return &Font{
		Face: f,
	}
}

func Inconsolata() *Font {
	return &Font{Face: inconsolataRegular}
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
