package txt

import (
	"github.com/ohzqq/pages"
)

func NewPaginator(txt string, opts ...WrapOpt) *pages.Pages[string] {
	wr := NewWrapper(opts...)
	return newPaginator(wr.WrapText(txt), wr.LinesPerPage())
}

func newPaginator(lines []string, linesPerPage int) *pages.Pages[string] {
	return pages.New(lines, pages.WithPerPage(linesPerPage))
}

func PaginateTextbox(txt string, w, h int, opts ...WrapOpt) *pages.Pages[string] {
	lines, pp := WrapTextbox(txt, w, h, opts...)
	return newPaginator(lines, pp)
}

func PaginateText(txt string, w int, maxLines int) *pages.Pages[string] {
	lines, pp := SimpleWrap(txt, w, maxLines)
	return newPaginator(lines, pp)
}
