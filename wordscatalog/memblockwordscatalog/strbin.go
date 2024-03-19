package memblockwordscatalog

import (
	"errors"
	"fmt"
	"main/datastructures"
	"main/wordscatalog"
	"sort"
)

type StrBin struct {
	// an intermediate type helping to create the MemBlockWordsCatalog
	commonStringsLength int
	strings             []string
	isSorted            bool
	isFinalized         bool // publishing an iterator finalizes the bin. we don't want concurrent modifications.
}

func (bin StrBin) AddString(str string) error {
	if bin.isFinalized {
		// the StrBin is thread safe only with regards to calls made after it was finalized.
		// no guaranties for concurrent iteration or modification
		return errors.New("Cannot call AddString after the StrBin has been finalized by publishing an iterator!.")
	}
	runes := []rune(str)
	if len(runes) == bin.commonStringsLength {
		bin.strings = append(bin.strings, str)
		bin.isSorted = false
		return nil
	} else {
		return errors.New(fmt.Sprintf("cannot add a string of length %v to a bin of strings with length %v", len(runes), bin.commonStringsLength))
	}
}

func (bin StrBin) sortBin() {
	// sorts the bin internally according to string comparison on his canonical form
	// the canonical form of a string is just the sorted string.

	sort.Slice(bin.strings, func(i, j int) bool {
		return wordscatalog.SortString(bin.strings[i]) < wordscatalog.SortString(bin.strings[j])
	})
	bin.isSorted = true
}

func (bin StrBin) getRequiredMemorySize() int {
	// assuming that each character in each string is UTF-8 ascii, thus takes 1 byte
	return bin.commonStringsLength * len(bin.strings)
}

func (bin StrBin) GetBytesIterator() datastructures.Iterator[byte] {
	// returns an iterator that trasverses the bytes in the strings according to their
	// order as if the strings were all sorted and laid out contiguously according to sorted order
	bin.isFinalized = true
	if !bin.isSorted {
		bin.sortBin()
	}
	currentStringIndex := 0
	currentByteIndex := 0

	hasNext := func() bool {
		return currentByteIndex < bin.commonStringsLength && currentStringIndex < len(bin.strings)
	}
	getNext := func() byte {
		currentStr := bin.strings[currentStringIndex]
		runes := []rune(currentStr)
		res := runes[currentByteIndex]
		currentByteIndex += 1
		if currentByteIndex >= bin.commonStringsLength {
			currentByteIndex = 0
			currentStringIndex += 1
		}
		return byte(res) // assuming the rune is an ascii UTF-8, the conversion will not truncate
	}

	getError := func() error { return nil }

	return datastructures.Iterator[byte]{HasNext: hasNext, GetNext: getNext, GetErr: getError}
}
