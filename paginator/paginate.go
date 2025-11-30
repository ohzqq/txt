package paginate

type Paginator[P any] struct {
	Pages    [][]P
	og       []P
	cur      int
	PageSize int
}

func New[T any](stuff []T) *Paginator[T] {
	return &Paginator[T]{og: stuff}
}
