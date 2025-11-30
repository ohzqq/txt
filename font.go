package main

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
)

var (
	GoMono      = gomono.TTF
	GoRegular   = goregular.TTF
	Inconsolata = inconsolata.Regular8x16
	DefaultFont = inconsolata.Regular8x16
)

func NewTTF(src []byte, opts *opentype.FaceOptions) (font.Face, error) {
	inc, err := opentype.Parse(src)
	if err != nil {
		return &opentype.Face{}, err
	}
	return opentype.NewFace(inc, opts)
}
