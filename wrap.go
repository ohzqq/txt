package txt

import (
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const elipsis = `...`

var (
	GoMono      = gomono.TTF
	GoRegular   = goregular.TTF
	Inconsolata = inconsolata.Regular8x16
)

func WrapToString(text string, w int) string {
	return strings.Join(WrapToSlice(text, w), "\n")
}

func WrapToSlice(str string, w int) []string {
	lines := split(str, w)
	return lines
}

func WrapAndChunk(str string, w, lh int) [][]string {
	return slice.Chunk(split(str, w), lh)
}

type WrapText struct {
	FontSize int
	Width    int
	DPI      int
	MaxLines int
	Font     font.Face
}

type Box struct {
	*WrapText
	Pages [][]string
}

func NewBox(text string, wr *WrapText) *Box {
	box := &Box{
		WrapText: wr,
	}
	box.Pages = box.Paginate(text)
	return box
}

func NewTextWrapper() *WrapText {
	box := &WrapText{
		Font:     NewFont(),
		FontSize: 16,
		DPI:      72,
		MaxLines: 3,
	}
	return box
}

func (box *WrapText) WrapText(text string) []string {
	return WrapLines(text, box.Font, box.Width)
}

func (box *WrapText) Paginate(text string) [][]string {
	lines := box.WrapText(text)
	if box.MaxLines <= 0 {
		return [][]string{lines}
	}
	return lo.Chunk(lines, box.MaxLines)
}

func (box *WrapText) SetMaxLines(m int) *WrapText {
	box.MaxLines = m
	return box
}

func (box *WrapText) SetWidth(w int) *WrapText {
	box.Width = w
	return box
}

func (box *WrapText) SetFont(face font.Face) *WrapText {
	box.Font = face
	return box
}

func (box *WrapText) WithTTF(src []byte) error {
	opts := &opentype.FaceOptions{
		Size:    float64(box.FontSize),
		DPI:     float64(box.DPI),
		Hinting: font.HintingNone,
	}
	f, err := NewTTFSrc(src, opts)
	if err != nil {
		return err
	}
	box.Font = f
	return nil
}

func WrapLines(str string, face font.Face, pxWidth int) []string {
	var f text.Frame
	f.SetFace(face)
	f.SetMaxWidth(fixed.I(pxWidth))
	c := f.NewCaret()
	c.WriteString(str)
	c.Close()
	txt := wrapBox(&f)
	return txt
}

func wrapBox(f *text.Frame) []string {
	txt := []string{}
	for p := f.FirstParagraph(); p != nil; p = p.Next(f) {
		for l := p.FirstLine(f); l != nil; l = l.Next(f) {
			line := []string{}
			for b := l.FirstBox(f); b != nil; b = b.Next(f) {
				line = append(line, string(b.TrimmedText(f)))
			}
			txt = append(txt, strings.Join(line, " "))
		}
	}
	return txt
}
