package txt

import (
	"github.com/ohzqq/pages"
)

func NewPaginator(lines []string, linesPerPage int) *pages.Pages[string] {
	return pages.New(lines, pages.WithPerPage(linesPerPage))
}
