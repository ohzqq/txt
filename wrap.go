package txt

import (
	"strings"

	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const elipsis = `...`

type TextWrapper interface {
	WrapText(text string) ([]string, int)
}

type Wrapper struct {
	FontSize   int
	Width      int
	DPI        int
	MaxLines   int
	Height     int
	LineHeight int
	Font       font.Face
	Simple     bool
}

type SimpleTextWrapper struct {
	Width int
}

func NewWrapper() *Wrapper {
	box := &Wrapper{
		Font:     DefaultFont,
		Width:    250,
		FontSize: 16,
		DPI:      72,
		MaxLines: 1,
	}
	return box
}

func (pg *Wrapper) LinesPerPage() int {
	if pg.Height > 0 {
		return pg.Height / pg.LineHeight
	}
	if pg.MaxLines > 1 {
		return pg.MaxLines
	}
	return 1
}

func (wr *Wrapper) SetFontSize(fs int) *Wrapper {
	wr.FontSize = fs
	return wr
}

func (wr *Wrapper) SetWidth(w int) *Wrapper {
	wr.Width = w
	return wr
}

func (wr *Wrapper) SetFont(face font.Face) *Wrapper {
	wr.Font = face
	return wr
}

func (pg *Wrapper) SetMaxLines(m int) *Wrapper {
	pg.MaxLines = m
	return pg
}

func (pg *Wrapper) SetLineHeight(m int) *Wrapper {
	pg.LineHeight = m
	return pg
}

func (pg *Wrapper) SetHeight(m int) *Wrapper {
	pg.Height = m
	return pg
}

func (wr *Wrapper) WithGoMono() *Wrapper {
	f, _ := NewTTF(GoMono, wr.opentypeOpts())
	return wr.SetFont(f)
}

func (wr *Wrapper) WithGoRegular() *Wrapper {
	f, _ := NewTTF(GoRegular, wr.opentypeOpts())
	return wr.SetFont(f)
}

func (wr *Wrapper) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(wr.FontSize),
		DPI:     float64(wr.DPI),
		Hinting: font.HintingNone,
	}
}

func (wr *Wrapper) WithTTF(src []byte) error {
	f, err := NewTTF(src, wr.opentypeOpts())
	if err != nil {
		return err
	}
	wr.Font = f
	return nil
}

func (wr *Wrapper) SimpleWrap(str string) []string {
	return split(str, wr.Width)
}

func (wr *Wrapper) WrapText(str string) []string {
	var f text.Frame
	f.SetFace(wr.Font)
	f.SetMaxWidth(fixed.I(wr.Width))
	c := f.NewCaret()
	c.WriteString(str)
	c.Close()
	lines, h := wrapBox(&f)
	wr.LineHeight = h
	return lines
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
