package txt

import (
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const elipsis = `...`

type TextWrapper interface {
	WrapText(text string) ([]string, int)
}

type WrapText struct {
	FontSize int
	Width    int
	DPI      int
	Font     font.Face
}

type Paginator struct {
	wr         *WrapText
	MaxLines   int
	LineHeight int
	Height     int
}

type Pages struct {
	Pages  [][]string
	cur    int
	Loop   bool
	Height int
	Width  int
}

func NewPages(pages [][]string) *Pages {
	return &Pages{
		Pages: pages,
		Loop:  true,
	}
}

func NewPaginator(wr *WrapText) *Paginator {
	box := &Paginator{
		wr:       wr,
		MaxLines: 5,
	}
	return box
}

func (pg *Paginator) SetMaxLines(m int) *Paginator {
	pg.MaxLines = m
	return pg
}

func (pg *Paginator) SetLineHeight(m int) *Paginator {
	pg.LineHeight = m
	return pg
}

func (pg *Paginator) SetHeight(m int) *Paginator {
	pg.Height = m
	return pg
}

func (pg *Paginator) Paginate(text string) *Pages {
	lines, height := pg.wr.WrapText(text)
	pg.LineHeight = height
	if pg.MaxLines <= 0 {
		return NewPages([][]string{lines})
	}
	if pg.Height > 0 {
		pg.MaxLines = pg.Height / height
	}
	return NewPages(lo.Chunk(lines, pg.MaxLines))
}

func (pg *Pages) SetLoop(loop bool) *Pages {
	pg.Loop = loop
	return pg
}

func (pg *Pages) CurrentLines() []string {
	return pg.Pages[pg.cur]
}

func (pg *Pages) CurrentPage() string {
	return strings.Join(pg.CurrentLines(), "\n")
}

func (pg *Pages) NextPage() string {
	pg.cur++
	if pg.cur >= len(pg.Pages) {
		if pg.Loop {
			pg.cur = 0
		} else {
			pg.cur = len(pg.Pages) - 1
		}
	}
	return pg.CurrentPage()
}

func (pg *Pages) PrevPage() string {
	pg.cur--
	if pg.cur < 0 {
		if pg.Loop {
			pg.cur = len(pg.Pages) - 1
		} else {
			pg.cur = 0
		}
	}
	return pg.CurrentPage()
}

func NewTextWrapper() *WrapText {
	box := &WrapText{
		Font:     DefaultFont,
		Width:    250,
		FontSize: 16,
		DPI:      72,
	}
	return box
}

func (box *WrapText) WrapText(text string) ([]string, int) {
	return WrapLines(text, box.Font, box.Width)
}

func (wr *WrapText) SetFontSize(fs int) *WrapText {
	wr.FontSize = fs
	return wr
}

func (box *WrapText) SetWidth(w int) *WrapText {
	box.Width = w
	return box
}

func (box *WrapText) SetFont(face font.Face) *WrapText {
	box.Font = face
	return box
}

func (box *WrapText) WithGoMono() *WrapText {
	f, _ := NewTTF(GoMono, box.opentypeOpts())
	return box.SetFont(f)
}

func (box *WrapText) WithGoRegular() *WrapText {
	f, _ := NewTTF(GoRegular, box.opentypeOpts())
	return box.SetFont(f)
}

func (wr *WrapText) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(wr.FontSize),
		DPI:     float64(wr.DPI),
		Hinting: font.HintingNone,
	}
}

func (box *WrapText) WithTTF(src []byte) error {
	f, err := NewTTF(src, box.opentypeOpts())
	if err != nil {
		return err
	}
	box.Font = f
	return nil
}

func WrapLines(str string, face font.Face, pxWidth int) ([]string, int) {
	var f text.Frame
	f.SetFace(face)
	f.SetMaxWidth(fixed.I(pxWidth))
	c := f.NewCaret()
	c.WriteString(str)
	c.Close()
	txt, height := wrapBox(&f)
	return txt, height
}

func wrapBox(f *text.Frame) ([]string, int) {
	txt := []string{}
	height := 0
	for p := f.FirstParagraph(); p != nil; p = p.Next(f) {
		for l := p.FirstLine(f); l != nil; l = l.Next(f) {
			line := []string{}
			height = l.Height(f)
			for b := l.FirstBox(f); b != nil; b = b.Next(f) {
				line = append(line, string(b.TrimmedText(f)))
			}
			txt = append(txt, strings.Join(line, " "))
		}
	}
	return txt, height
}
