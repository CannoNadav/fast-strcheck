package memblockwordscatalog

import (
	"fmt"
	"main/wordscatalog"
	"sort"
	"strings"
	"unsafe"
)

type WordsBlock struct {
	memBlock   []byte // view of the byte slice holding the strings of this block
	numStrings int    // amount of strings in the block.
	strLength  int    // the length, shared by all strings in the block
}

func (wordsBlock WordsBlock) getNumEquivalenceGroups() int {
	numGroups := 0
	for index := 0; index < wordsBlock.numStrings; {
		str := wordsBlock.getIthString(index)
		equivGroup := wordsBlock.FindEquivalentWords(str)
		numGroups += 1
		index += len(equivGroup) + 1 // the +1 is for the str itself, that isn't extracted as part of his own equivalence group, but is still laid out in memory with the rest of
	}
	//print(fmt.Sprintf("wordsBlock[strlength=%v]: %v equiv groups\n", wordsBlock.strLength, numGroups))
	return numGroups

}

func (wordsBlock WordsBlock) getMemorySize() int {
	return len(wordsBlock.memBlock) + int(unsafe.Sizeof(wordsBlock))
}

func (wordsBlock WordsBlock) getNumStringsInBlock() int {
	return wordsBlock.numStrings
}

func (wordsBlock WordsBlock) getIthString(idx int) string {
	// idx must be a valid index for one of the values in the WordsBlock
	//println(fmt.Sprintf("asked to retrieve str num: %v out of %v", idx, wordsBlock.numStrings))
	if idx >= wordsBlock.numStrings {
		panic(fmt.Sprintf("invalid index %v for a WordsBlock with %v strings", idx, wordsBlock.numStrings))
	}
	bytes := wordsBlock.memBlock[idx*wordsBlock.strLength : (idx+1)*wordsBlock.strLength]
	// since the memory is internal and never exposed or modified, I can avoid
	// the cost of frequent string copying
	//firstByteInStrIndex := idx*wordsBlock.strLength
	//print(fmt.Sprintf("first byte in string in position %v out of %v",firstByteInStrIndex, len() )
	//firstByteInString := bytes[idx*wordsBlock.strLength]
	return unsafe.String(&bytes[0], wordsBlock.strLength)
	//return *(*string)(unsafe.Pointer(&bytes))
}

func (wordsBlock WordsBlock) FindEquivalentWords(word string) []string {
	var ithString string
	wordSortedLexicographically := wordscatalog.SortString(word)
	isEquivWord := func(i int) int {
		ithString = wordsBlock.getIthString(i)
		ithStringSortedLexicographically := wordscatalog.SortString(ithString)
		return strings.Compare(wordSortedLexicographically, ithStringSortedLexicographically)
	}
	firstIndex, isFound := sort.Find(wordsBlock.numStrings, isEquivWord)
	result := []string(nil)
	if !isFound {
		return result
	}
	// now collect the equivalence group from the underlying "array",
	// equivalent words are stored one after the other in the array

	for currWordIdx := firstIndex; currWordIdx < wordsBlock.numStrings && isEquivWord(currWordIdx) == 0; currWordIdx++ {
		var currentWord string
		currentWord = wordsBlock.getIthString(currWordIdx)
		if currentWord != word { // The query word should be excluded from the equivalent words result list
			result = append(result, currentWord)
		}
	}
	return result
}
