package txt

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
)

func NewFont() *basicfont.Face {
	font := inconsolata.Regular8x16
	return font
}

func NewTTF(file string, size int) (font.Face, error) {
	ff, err := os.ReadFile(file)
	if err != nil {
		return &opentype.Face{}, err
	}

	return NewTTFSrc(ff, nil)
}

func NewTTFSrc(src []byte, opts *opentype.FaceOptions) (font.Face, error) {
	inc, err := opentype.Parse(src)
	if err != nil {
		return &opentype.Face{}, err
	}
	return opentype.NewFace(inc, opts)
}
