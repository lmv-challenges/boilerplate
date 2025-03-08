package boilerplate

import (
	"fmt"
	"strings"
)

type SliceToSliceFunc[T any] func(in []T) (out []T)

type SliceScanner[T any] func() (out []T, err error)

func WithLengthAndSliceScanner[T any]() ([]T, error) {
	var n uint
	_, err := fmt.Scan(&n)
	if err != nil {
		return nil, fmt.Errorf("failed to scan n: %w", err)
	}
	nums := make([]T, n)
	for i := range n {
		_, err = fmt.Scan(&nums[i])
		if err != nil {
			var zero T
			return nil, fmt.Errorf("failed to read %T i=%d: %w", zero, i, err)
		}
	}
	return nums, nil
}

type SlicePrinter[T any] func(out []T)

func WithSlicePrinter[T any](out []T) {
	strOut := make([]string, len(out))
	for i, val := range out {
		strOut[i] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strOut, " "))
}

func ReadIntAndSliceWriteSlice[T any](f SliceToSliceFunc[T], r SliceScanner[T], w SlicePrinter[T]) error {
	in, err := r()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	out := f(in)
	w(out)
	return nil
}
