package memblockwordscatalog

import (
	"errors"
	"fmt"
	"main/datastructures"
	"main/wordscatalog"
	"sort"
	"unsafe"
)

// since cap=0 theEmptySlice will always stay empty. if the client adds to it, then a new underlying
// slice will be allocated for that client, as the slice grew.
var theEmptySlice []string

type MemBlockWordsCatalog struct {
	stringBlocksMap map[int]WordsBlock
}

func (memBlockCatalog MemBlockWordsCatalog) GetNumEquivalenceClasses() int {
	numGroups := 0
	for _, block := range memBlockCatalog.stringBlocksMap {
		numGroups += block.getNumEquivalenceGroups()
	}
	return numGroups
}

func (memBlockCatalog MemBlockWordsCatalog) GetMemorySize() int {
	mem := 0
	for _, wordBlock := range memBlockCatalog.stringBlocksMap {
		mem += wordBlock.getMemorySize()
	}
	return mem + int(unsafe.Sizeof(memBlockCatalog))
}

//func (catalog MemBlockWordsCatalog) BankWordsVoucher(result []string, voucher WordsCatalogVoucher) []string {
//	// gets the words represented in the voucher and adds them to the slice result.
//	// allows the caller to reuse slices for subsequent calls.
//	for i := 0; i<voucher.numStrings; i++{
//		result = append(result, voucher.strBlock.getIthString(voucher.indexInBlock + i))
//	}
//	return result
//}

//type WordsCatalogVoucher struct {
//	// never instantiate directly!! constructor should be private
//	strBlock WordsBlock
//	indexInBlock int
//	numStrings int
//}
//
//func (voucher WordsCatalogVoucher) GetNumStringsInVoucher() int {
//	return voucher.numStrings
//}

func (memBlockCatalog MemBlockWordsCatalog) FindEquivalentWords(word string) []string {
	strLen := len([]rune(word))
	strBlock, exists := memBlockCatalog.stringBlocksMap[strLen]
	if !exists {
		return theEmptySlice
	}

	res := strBlock.FindEquivalentWords(word)
	if res == nil || len(res) == 0 {
		return theEmptySlice
	}
	return res
}

func (memBlockWordsCatalog MemBlockWordsCatalog) CountWordsInCatalog() int {
	wordsCount := 0
	for _, strBlock := range memBlockWordsCatalog.stringBlocksMap {
		wordsCount += strBlock.getNumStringsInBlock()
	}
	return wordsCount
}

func isASCII(r rune) bool {
	return r >= 0 && r <= 127
}

func validateEncodingAssumptions(str string) error {
	// checks if the string str contains a non-ascii utf-8 character.  returns an error if true.
	// the datastructure is built under the assumption that all characters take
	// 1 byte (ascii characters take 1 byte in UTF-8 encoding).

	runes := []rune(str)
	for _, r := range runes {
		if !isASCII(r) {
			return errors.New(fmt.Sprintf("The string %v contained a non ascii character %v. expecting only ascii characters", str, r))
		}
	}
	return nil
}

func groupStringsByLength(strings []string) map[int][]string {
	stringsBinnedByLength := map[int][]string{}
	for _, str := range strings {
		runes := []rune(str)
		strLength := len(runes)
		bin, exists := stringsBinnedByLength[strLength]
		if !exists {
			bin = []string{}
		}
		bin = append(bin, str)
		stringsBinnedByLength[strLength] = bin
	}

	return stringsBinnedByLength

}

func NewMemBlockWordsCatalog(wordsIter datastructures.Iterator[string]) (wordscatalog.WordsCatalog, error) {
	theEmptySlice = make([]string, 0)
	allWordsInCatalog := []string{}
	for wordsIter.HasNext() {
		currentWord := wordsIter.GetNext()
		allWordsInCatalog = append(allWordsInCatalog, currentWord)
		err := validateEncodingAssumptions(currentWord)
		if err != nil {
			return MemBlockWordsCatalog{}, err
		}
	}
	if len(allWordsInCatalog) == 0 {
		return MemBlockWordsCatalog{}, errors.New("0 strings were sent as input to build the catalog . is this what you wanted to do?")
	}

	allStrBins := []StrBin{}
	for strLength, stringSlice := range groupStringsByLength(allWordsInCatalog) {
		bin := StrBin{commonStringsLength: strLength, strings: stringSlice}
		allStrBins = append(allStrBins, bin)
	}
	// sorts the slice. has no effect on correctness, mostly for easier inspection
	sort.Slice(allStrBins, func(i, j int) bool {
		return allStrBins[i].commonStringsLength < allStrBins[j].commonStringsLength
	})

	// count how many space we need to store all the strings
	totalBytesNeeded := 0
	for _, bin := range allStrBins {
		totalBytesNeeded += bin.getRequiredMemorySize()
	}

	// now allocate one contiguous memory block and put all the strings into it
	wordsBlocksMap := map[int]WordsBlock{}
	memoryBlock := make([]byte, totalBytesNeeded)
	currentBlockIndex := 0
	for _, strBin := range allStrBins {
		byteIterator := strBin.GetBytesIterator() // this action finalizes the StrBin and sorts it internally to have proper memory layout for later search before writing it to the block
		startIndex := currentBlockIndex
		for byteIterator.HasNext() {
			b := byteIterator.GetNext()
			memoryBlock[currentBlockIndex] = b
			currentBlockIndex += 1
		}
		endIndex := currentBlockIndex
		wordsBlock := WordsBlock{memBlock: memoryBlock[startIndex:endIndex], numStrings: len(strBin.strings), strLength: strBin.commonStringsLength}
		wordsBlocksMap[strBin.commonStringsLength] = wordsBlock
	}

	return MemBlockWordsCatalog{stringBlocksMap: wordsBlocksMap}, nil

}
