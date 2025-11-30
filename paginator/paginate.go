package paginate

import "github.com/samber/lo"

type Pages[P any] struct {
	*Paginator
	pages []P
}

func New[T any](stuff []T, opts ...Option) *Pages[T] {
	pagy := NewPaginator(opts...)
	pagy.SetTotalPages(len(stuff))
	return &Pages[T]{
		pages:     stuff,
		Paginator: &pagy,
	}
}

func (p Pages[T]) Current() []T {
	s, e := p.GetSliceBounds(len(p.pages))
	return p.pages[s:e]
}

func (p Pages[T]) AllPages() [][]T {
	return lo.Chunk(p.pages, p.PerPage)
}

func (p Pages[T]) Next() []T {
	n := p.nextPageNum()
	return p.AllPages()[n]
}

func (p Pages[T]) Prev() []T {
	n := p.prevPageNum()
	return p.AllPages()[n]
}
