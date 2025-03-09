package boilerplate

import (
	"bufio"
	"fmt"
	"strings"
)

type ProcessFunc[T any, U any] func(in []T) (out []U)

type SliceScanner[U any] func() (out []U, err error)

func WithLengthAndSliceScanner[X any](reader *bufio.Reader) ([]X, error) {
	var n uint
	_, err := fmt.Fscan(reader, &n)
	if err != nil {
		return nil, fmt.Errorf("failed to scan n: %w", err)
	}
	elems := make([]X, n)
	for i := range n {
		_, err = fmt.Fscan(reader, &elems[i])
		if err != nil {
			var zero X
			return nil, fmt.Errorf("failed to read %T i=%d: %w", zero, i, err)
		}
	}
	return elems, nil
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

func WithTestCases(f func() error, reader *bufio.Reader) error {
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
