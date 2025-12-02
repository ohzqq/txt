package txt

import (
	"strings"

	"github.com/ohzqq/pages"
)

type Pages struct {
	*pages.Pages[string]
}

func NewPaginator(lines []string, linesPerPage int) *Pages {
	return &Pages{
		Pages: pages.New(lines, pages.WithPerPage(linesPerPage)),
	}
}

func (p *Pages) JoinCurrentWithSpace() string {
	return strings.TrimSpace(strings.Join(p.Current(), " "))
}

func (p *Pages) JoinCurrentWithNewLine() string {
	return strings.TrimSpace(strings.Join(p.Current(), "\n"))
}
