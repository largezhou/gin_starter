package v

// P 返回一个值的指针
func P[T any](v T) *T {
	return &v
}

// V 返回一个指针的底层值，指针为 nil 则返回 0值
func V[T any](p *T) T {
	if p == nil {
		var t T
		return t
	}

	return *p
}
