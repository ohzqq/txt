// Package paginator provides a Bubble Tea package for calculating pagination
// and rendering pagination info. Note that this package does not render actual
// pages: it's purely for handling keystrokes related to pagination, and
// rendering pagination status.

//MIT License

//Copyright (c) 2020-2025 Charmbracelet, Inc

//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package paginate

// Paginator is the Bubble Tea model for this user interface.

type Paginator struct {
	// Page is the current page number.
	Page int
	// PerPage is the number of items per page.
	PerPage int
	// TotalPages is the total number of pages.
	TotalPages int
	// Loop determines whether or not to go to the first page, after the last
	Loop bool
}

// SetTotalPages is a helper function for calculating the total number of pages
// from a given number of items. Its use is optional since this pager can be
// used for other things beyond navigating sets. Note that it both returns the
// number of total pages and alters the model.
func (m *Paginator) SetTotalPages(items int) int {
	if items < 1 {
		return m.TotalPages
	}
	n := items / m.PerPage
	if items%m.PerPage > 0 {
		n++
	}
	m.TotalPages = n
	return n
}

// ItemsOnPage is a helper function for returning the number of items on the
// current page given the total number of items passed as an argument.
func (m Paginator) ItemsOnPage(totalItems int) int {
	if totalItems < 1 {
		return 0
	}
	start, end := m.GetSliceBounds(totalItems)
	return end - start
}

// GetSliceBounds is a helper function for paginating slices. Pass the length
// of the slice you're rendering and you'll receive the start and end bounds
// corresponding to the pagination. For example:
//
//	bunchOfStuff := []stuff{...}
//	start, end := model.GetSliceBounds(len(bunchOfStuff))
//	sliceToRender := bunchOfStuff[start:end]
func (m *Paginator) GetSliceBounds(length int) (start int, end int) {
	start = m.Page * m.PerPage
	end = min(m.Page*m.PerPage+m.PerPage, length)
	return start, end
}

// PrevPage is a helper function for navigating one page backward. It will not
// page beyond the first page (i.e. page 0).
func (m *Paginator) PrevPage() {
	m.Page = m.PrevPageNum()
}

func (m *Paginator) PrevPageNum() int {
	prev := m.Page
	prev--
	if prev < 0 {
		if m.Loop {
			prev = m.TotalPages - 1
		} else {
			prev = 0
		}
	}
	return prev
}

// NextPage is a helper function for navigating one page forward. It will not
// page beyond the last page (i.e. totalPages - 1).
func (m *Paginator) NextPage() {
	m.Page = m.NextPageNum()
}

func (m *Paginator) NextPageNum() int {
	next := m.Page
	if !m.OnLastPage() {
		next++
	} else if m.Loop {
		next = 0
	}
	return next
}

// OnLastPage returns whether or not we're on the last page.
func (m Paginator) OnLastPage() bool {
	return m.Page == m.TotalPages-1
}

// OnFirstPage returns whether or not we're on the first page.
func (m Paginator) OnFirstPage() bool {
	return m.Page == 0
}

// Option is used to set options in New.
type Option func(*Paginator)

// NewPaginator creates a new model with defaults.
func NewPaginator(opts ...Option) Paginator {
	m := Paginator{
		Page:       0,
		PerPage:    1,
		TotalPages: 1,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// WithLooping sets pagination to loop.
func WithLooping(loop bool) Option {
	return func(m *Paginator) {
		m.Loop = loop
	}
}

// WithTotalPages sets the total pages.
func WithTotalPages(totalPages int) Option {
	return func(m *Paginator) {
		m.TotalPages = totalPages
	}
}

// WithPerPage sets the total pages.
func WithPerPage(perPage int) Option {
	return func(m *Paginator) {
		m.PerPage = perPage
	}
}
