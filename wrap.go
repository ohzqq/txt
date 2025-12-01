package txt

import (
	"strings"

	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/math/fixed"
)

type Wrapper struct {
	MaxLines   int
	LineHeight int
	Height     int
	Simple     bool
}

func NewWrapper() *Wrapper {
	wr := &Wrapper{
		MaxLines: 1,
	}
	return wr
}

//func SimpleWrap(txt string, w, maxLines int) ([]string, int) {
//  wr := &Wrapper{
//    Font: &Font{
//      Width: 250,
//    },
//    MaxLines: maxLines,
//    Simple:   true,
//  }
//  return wr.wrapText(txt), wr.LinesPerPage()
//}

//func WrapFont(txt string, opts ...WrapOpt) ([]string, int) {
//  wr := NewWrapper(opts...)
//  if wr.Font == nil {
//    wr.Font = Inconsolata()
//  }
//  return wr.wrapText(txt), wr.LinesPerPage()
//}

//func WrapTextbox(txt string, w, h int, opts ...WrapOpt) ([]string, int) {
//  wr := NewWrapper(opts...)
//  if wr.Font == nil {
//    wr.Font = Inconsolata()
//  }
//  wr.Width = w
//  wr.Height = h
//  return wr.wrapText(txt), wr.LinesPerPage()
//}

func (wr *Wrapper) WrapBytes(byte []byte, font *Font) []string {
	return wr.WrapText(string(byte), font)
}

func (wr *Wrapper) WrapText(str string, font *Font) []string {
	if wr.Simple == true {
		return simpleWrap(str, font.Width)
	}
	var f text.Frame
	f.SetFace(font.Face)
	f.SetMaxWidth(fixed.I(font.Width))
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
