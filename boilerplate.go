package boilerplate

import (
	"fmt"
	"io"
	"strings"
)

type ProcessFunc[T any, U any] func(in []T) (out []U)

type SliceScanner[U any] func() (out []U, err error)

// WithLengthAndSliceScanner returns a closure that reads an integer n from reader followed by n elements of type U.
// The closure returns a slice of U containing the n elements.
func WithLengthAndSliceScanner[U any](reader io.Reader) SliceScanner[U] {
	return func() ([]U, error) {
		var n uint
		_, err := fmt.Fscan(reader, &n)
		if err != nil {
			return nil, fmt.Errorf("failed to scan n: %w", err)
		}
		elems := make([]U, n)
		for i := range n {
			_, err = fmt.Fscan(reader, &elems[i])
			if err != nil {
				var zero U
				return nil, fmt.Errorf("failed to read element type=%T i=%d: %w", zero, i, err)
			}
		}
		return elems, nil
	}
}

// WithSliceScanner returns a closure that reads n elements of type U and returns them in a slice.
func WithSliceScanner[U any](reader io.Reader) SliceScanner[U] {
	return func() ([]U, error) {
		var elems []U
		for i := 0; ; i++ {
			var elem U
			_, err := fmt.Fscan(reader, &elem)
			if err == io.EOF {
				return elems, nil
			}
			if err != nil {
				var zero U
				return nil, fmt.Errorf("failed to read element type=%T i=%d: %w", zero, i, err)
			}
			elems = append(elems, elem)
		}
	}
}

type SlicePrinter[U any] func(out []U)

func WithSlicePrinter[U any](out []U) {
	strOut := make([]string, len(out))
	for i, val := range out {
		strOut[i] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strOut, " "))
}

func ScanProcessPrint[T any, U any](f ProcessFunc[T, U], r SliceScanner[T], w SlicePrinter[U]) error {
	in, err := r()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	out := f(in)
	w(out)
	return nil
}

func WithTestCases(f func() error, reader io.Reader) error {
	var n uint
	_, err := fmt.Fscan(reader, &n)
	if err != nil {
		return fmt.Errorf("failed to scan n test cases: %w", err)
	}
	for i := range n {
		if err = f(); err != nil {
			return fmt.Errorf("failed to process test case i=%d: %w", i, err)
		}
	}
	return nil
}
