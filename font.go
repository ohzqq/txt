package txt

import (
	"golang.org/x/image/font/basicfont"
)

func NewFont(size int) *basicfont.Face {
	font := basicfont.Face7x13
	font.Width = size
	return font
}
