package txt

import (
	"strings"

	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Wrapper struct {
	MaxLines   int
	LineHeight int
	Height     int
	Simple     bool
}

func NewWrapper() *Wrapper {
	return &Wrapper{
		MaxLines:   1,
		LineHeight: 1,
	}
}

func (wr *Wrapper) WrapBytes(byte []byte, font *Font) []string {
	return wr.WrapText(string(byte), font)
}

func WrapText(str string, face font.Face, w int) ([]string, int) {
	var f text.Frame
	f.SetFace(face)
	f.SetMaxWidth(fixed.I(w))
	c := f.NewCaret()
	c.WriteString(str)
	c.Close()
	return wrapBox(&f)
}

func (wr *Wrapper) WrapText(str string, fnt *Font) []string {
	face, _ := opentype.NewFace(fnt.font, fnt.opentypeOpts())
	if wr.Simple == true {
		return simpleWrap(str, fnt.Width)
	}
	lines, h := WrapText(str, face, fnt.Width)
	wr.LineHeight = h
	fnt.LineHeight = h
	fnt.linesPerPage = calculateLinesPerPage(fnt.Height, fnt.LineHeight)
	return lines
}

func calculateLinesPerPage(height, lineHeight int) int {
	if height > 0 {
		return height / lineHeight
	}
	return 1
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

func (pg *Wrapper) SetMaxLines(m int) *Wrapper {
	pg.MaxLines = m
	return pg
}

func (pg *Wrapper) SetLineHeight(m int) *Wrapper {
	pg.LineHeight = m
	return pg
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
