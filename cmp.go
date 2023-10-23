package cmp

import "cmp"

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
//
// Note that floating-point types may contain NaN ("not-a-number") values.
// An operator such as == or < will always report false when
// comparing a NaN value with any other value, NaN or not.
// See the [Natural] function for a consistent way to compare NaN values.
type Ordered = cmp.Ordered

// Comparator is a function that compares two values. It returns
// -1 if a < b, 0 if a == b, 1 if a > b
type Comparator[T any] func(a, b T) int

func By[T any, R Ordered](f func(T) R) Comparator[T] {
	return func(a, b T) int {
		return cmp.Compare(f(a), f(b))
	}
}

// Natural comparator returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Natural[T Ordered]() Comparator[T] {
	return cmp.Compare
}

// Less reports whether x is less than y.
func (c Comparator[T]) Less(a, b T) bool {
	return c(a, b) < 0
}

func (c Comparator[T]) Equal(a, b T) bool {
	return c(a, b) == 0
}

func (c Comparator[T]) Greater(a, b T) bool {
	return c(a, b) > 0
}

func (c Comparator[T]) Then(c2 Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		if res := c(a, b); res != 0 {
			return res
		}
		return c2(a, b)
	}
}

func (c Comparator[T]) Reversed() Comparator[T] {
	return func(a, b T) int {
		return -c(a, b)
	}
}

func (c Comparator[T]) WithBottom(t T) Comparator[T] {
	return func(a, b T) int {
		if c(a, t) == 0 {
			return -1
		}
		if c(b, t) == 0 {
			return 1
		}
		return c(a, b)
	}
}

func (c Comparator[T]) WithTop(t T) Comparator[T] {
	return func(a, b T) int {
		if c(a, t) == 0 {
			return 1
		}
		if c(b, t) == 0 {
			return -1
		}
		return c(a, b)
	}
}

// cant be method
func Ptr[T any](c Comparator[T]) Comparator[*T] {
	return func(a, b *T) int {
		switch {
		case a == nil && b == nil:
			return 0
		case a == nil:
			return -1
		case b == nil:
			return 1
		default:
			return c(*a, *b)
		}
	}
}

func (c Comparator[T]) Max(value T, values ...T) T {
	max := value
	for _, v := range values {
		if c(v, max) > 0 {
			max = v
		}
	}
	return max
}

func (c Comparator[T]) Min(value T, values ...T) T {
	min := values[0]
	for _, v := range values {
		if c(v, min) < 0 {
			min = v
		}
	}
	return min
}
