package datastructures

import (
	"bufio"
	"os"
)

type Iterator[T any] struct {
	HasNext func() bool
	GetNext func() T
	GetErr  func() error
}

func NewSliceIterator[T any](slc []T) Iterator[T] {
	index := 0
	return Iterator[T]{
		HasNext: func() bool {
			return index < len(slc)
		},
		GetNext: func() T {
			if index < len(slc) {
				nextVal := slc[index]
				index++
				return nextVal
			}
			return *new(T) // Return zero value for T
		},
		GetErr: func() error {
			return nil
		},
	}
}

func NewFileLineIterator(catalogFile *os.File) Iterator[string] {
	scanner := bufio.NewScanner(catalogFile)
	lineIterator := Iterator[string]{HasNext: scanner.Scan, GetNext: scanner.Text, GetErr: scanner.Err}
	return lineIterator
}
