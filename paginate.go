package txt

import (
	"github.com/ohzqq/pages"
)

func NewPaginator(txt string, opts ...FontOpt) *pages.Pages[string] {
	//wr := NewWrapper(opts...)
	wr := NewFont(opts...)
	return newPaginator(wr.WrapText(txt), wr.wrapper.LinesPerPage())
}

func newPaginator(lines []string, linesPerPage int) *pages.Pages[string] {
	return pages.New(lines, pages.WithPerPage(linesPerPage))
}
