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

// Current returns the current page slice.
func (p Pages[T]) Current() []T {
	s, e := p.GetSliceBounds(len(p.pages))
	return p.pages[s:e]
}

// AllPages returns all pages as a slice of slices.
func (p Pages[T]) AllPages() [][]T {
	return lo.Chunk(p.pages, p.PerPage)
}

// Next returns the next page, without changing the current page number. Use is
// for pre-rendering the next page.
func (p Pages[T]) Next() []T {
	return p.AllPages()[p.NextPageNum()]
}

// Prev returns the prev page, without changing the current page number. Use is
// for pre-rendering the prev page.
func (p Pages[T]) Prev() []T {
	return p.AllPages()[p.PrevPageNum()]
}
