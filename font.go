package txt

import (
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
)

var (
	goMono             = gomono.TTF
	goRegular          = goregular.TTF
	inconsolataRegular = inconsolata.Regular8x16
)

type Font struct {
	fg color.Color
	bg color.Color
	font.Face
}

func NewFont(src []byte, opts *opentype.FaceOptions) (*Font, error) {
	f, err := NewTTF(src, opts)
	if err != nil {
		return nil, err
	}
	return &Font{
		Face: f,
	}, nil
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

func NewTTF(src []byte, opts *opentype.FaceOptions) (font.Face, error) {
	inc, err := opentype.Parse(src)
	if err != nil {
		return &opentype.Face{}, err
	}
	return opentype.NewFace(inc, opts)
}
