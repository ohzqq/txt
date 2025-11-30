package txt

import (
	"strings"

	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const elipsis = `...`

type WrapOpt func(w *Wrapper)

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

func NewWrapper(opts ...WrapOpt) *Wrapper {
	wr := &Wrapper{
		Width:    250,
		FontSize: 16,
		DPI:      72,
		MaxLines: 1,
	}
	for _, opt := range opts {
		opt(wr)
	}
	return wr
}

func SimpleWrap(txt string, w, maxLines int) ([]string, int) {
	wr := &Wrapper{
		Width:    w,
		MaxLines: maxLines,
		Simple:   true,
	}
	return wr.WrapText(txt), wr.LinesPerPage()
}

func WrapFont(txt string, opts ...WrapOpt) ([]string, int) {
	wr := NewWrapper(opts...)
	if wr.Font == nil {
		wr.Font = Inconsolata
	}
	return wr.WrapText(txt), wr.LinesPerPage()
}

func WrapTextbox(txt string, w, h int, opts ...WrapOpt) ([]string, int) {
	wr := NewWrapper(opts...)
	if wr.Font == nil {
		wr.Font = Inconsolata
	}
	wr.Width = w
	wr.Height = h
	return wr.WrapText(txt), wr.LinesPerPage()
}

func (wr *Wrapper) SimpleWrap(str string) []string {
	wr.Simple = true
	return wr.WrapText(str)
}

func (wr *Wrapper) WrapText(str string) []string {
	if wr.Simple == true {
		return simpleWrap(str, wr.Width)
	}
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

func (wr *Wrapper) SetSize(w, h int) *Wrapper {
	wr.Width = w
	wr.Height = h
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

func (wr *Wrapper) WithTTF(src []byte, fs int) error {
	wr.SetFontSize(fs)
	f, err := NewTTF(src, wr.opentypeOpts())
	if err != nil {
		return err
	}
	wr.SetFont(f)
	return nil
}

func WithSize(w, h int) WrapOpt {
	return func(wr *Wrapper) {
		wr.Width = w
		wr.Height = h
	}
}

func WithMaxLines(maxLines int) WrapOpt {
	return func(wr *Wrapper) {
		wr.MaxLines = maxLines
	}
}

func WithFont(face font.Face) WrapOpt {
	return func(wr *Wrapper) {
		wr.SetFont(face)
	}
}

func WithGoMono(fs int) WrapOpt {
	return func(wr *Wrapper) {
		wr.WithTTF(GoMono, fs)
	}
}

func WithGoRegular(fs int) WrapOpt {
	return func(wr *Wrapper) {
		wr.WithTTF(GoRegular, fs)
	}
}

func (wr *Wrapper) opentypeOpts() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(wr.FontSize),
		DPI:     float64(wr.DPI),
		Hinting: font.HintingNone,
	}
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

// borrowed from https://gist.github.com/AmrSaber/2468f546fb67dc31576a14e1209870e6
func simpleWrap(str string, size int) []string {
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
