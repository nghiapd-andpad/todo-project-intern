package cast

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T {
	return &v
}

// Value returns the value of p.
// It returns the zero value if p is nil.
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
