package pagination

type Page[T any] struct {
	Items      []T
	TotalCount int64
	Offset     int
	Limit      int
	IsFirst    bool
	IsLast     bool
	HasNext    bool
	HasPrev    bool
}

func New[T any](items []T, total int64, offset, limit int) *Page[T] {
	if items == nil {
		items = []T{}
	}
	isFirst := offset == 0
	isLast := int64(offset+len(items)) >= total

	return &Page[T]{
		Items:      items,
		TotalCount: total,
		Offset:     offset,
		Limit:      limit,
		IsFirst:    isFirst,
		IsLast:     isLast,
		HasNext:    !isLast,
		HasPrev:    !isFirst,
	}
}
