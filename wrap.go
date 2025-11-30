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

type VariableTextWrapper struct {
	FontSize int
	Width    int
	DPI      int
	Font     font.Face
}

type Paginator struct {
	wr         TextWrapper
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

func NewPaginator(wr TextWrapper) *Paginator {
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

func NewTextWrapper() *VariableTextWrapper {
	box := &VariableTextWrapper{
		Font:     DefaultFont,
		Width:    250,
		FontSize: 16,
		DPI:      72,
	}
	return box
}

func (wr *VariableTextWrapper) SetFontSize(fs int) *VariableTextWrapper {
	wr.FontSize = fs
	return wr
}

func (wr *VariableTextWrapper) SetWidth(w int) *VariableTextWrapper {
	wr.Width = w
	return wr
}

func (wr *VariableTextWrapper) SetFont(face font.Face) *VariableTextWrapper {
	wr.Font = face
	return wr
}

func (wr *VariableTextWrapper) WithGoMono() *VariableTextWrapper {
	f, _ := NewTTF(GoMono, wr.opentypeOpts())
	return wr.SetFont(f)
}

func (wr *VariableTextWrapper) WithGoRegular() *VariableTextWrapper {
	f, _ := NewTTF(GoRegular, wr.opentypeOpts())
	return wr.SetFont(f)
}

func (wr *VariableTextWrapper) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(wr.FontSize),
		DPI:     float64(wr.DPI),
		Hinting: font.HintingNone,
	}
}

func (wr *VariableTextWrapper) WithTTF(src []byte) error {
	f, err := NewTTF(src, wr.opentypeOpts())
	if err != nil {
		return err
	}
	wr.Font = f
	return nil
}

func (wr *VariableTextWrapper) WrapText(str string) ([]string, int) {
	var f text.Frame
	f.SetFace(wr.Font)
	f.SetMaxWidth(fixed.I(wr.Width))
	c := f.NewCaret()
	c.WriteString(str)
	c.Close()
	return wrapBox(&f)
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

type SimpleTextWrapper struct {
	Width int
}

func NewSimpleTextWrapper(w int) *SimpleTextWrapper {
	return &SimpleTextWrapper{Width: w}
}

func (wr *SimpleTextWrapper) WrapText(str string) ([]string, int) {
	return split(str, wr.Width), 1
}

// borrowed from https://gist.github.com/AmrSaber/2468f546fb67dc31576a14e1209870e6
func split(str string, size int) []string {
	if size < 1 {
		return []string{}
	}
	str = strings.TrimSpace(str)
	start := 0
	chunks := make([]string, 0, len(str)/size)
	for start < len(str) {
		end := start + size
		if end >= len(str) {
			chunks = append(chunks, str[start:])
			break
		}
		// If the next character is a delimiter, take it
		// this is to avoid adding "-" at the end when the next character is a delimiter anyway
		next := str[end]
		if next == ' ' || next == '\n' {
			end++
		}
		chunk := str[start:end]
		cutWord := false
		// Try to find a new line within the limit
		length := strings.LastIndex(chunk, "\n")
		// If no new line found, try to find a space
		if length == -1 {
			length = strings.LastIndex(chunk, " ")
			// If no space found, then just split the text
			if length == -1 {
				length = len(chunk) - 1 // leave space for "-" character that will be appended
				cutWord = true
			}
		}
		chunk = chunk[:length]
		start += length
		if cutWord {
			chunk += "-"
		} else {
			// Ignore the space that we stopped at
			start++
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}
