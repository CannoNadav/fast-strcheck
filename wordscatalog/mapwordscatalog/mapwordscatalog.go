package mapwordscatalog

import (
	"main/datastructures"
	"main/wordscatalog"
)

var theEmptySlice []string

type MapBackedWordsCatalog struct {
	maxWordLength      int
	EquivalenceClasses map[string][]string
}

func (catalog MapBackedWordsCatalog) GetNumEquivalenceClasses() int {
	return len(catalog.EquivalenceClasses)
}

func (catalog MapBackedWordsCatalog) FindEquivalentWords(word string) []string {
	// shortcutting the search on word length, to prevent an attacker from sending very very long inputs
	// that might bog down the server when the string is sorted
	if len([]rune(word)) <= catalog.maxWordLength {

		cannonicalRepresentation := wordscatalog.SortString(word)
		if equivalentWords, keyExists := catalog.EquivalenceClasses[cannonicalRepresentation]; keyExists {
			// removing the parameter word, as it should not be part of the result
			equivClsWithoutSearchWord := []string{}
			for _, w := range equivalentWords {
				if w != word {
					equivClsWithoutSearchWord = append(equivClsWithoutSearchWord, w)
				}
			}

			return equivClsWithoutSearchWord
		}
	}
	return theEmptySlice
}

func NewMapBackedWordsCatalog(wordsIter datastructures.Iterator[string]) (wordscatalog.WordsCatalog, error) {
	theEmptySlice = make([]string, 0)
	maxWordLength := -1
	equivalenceClasses := make(map[string][]string)
	for wordsIter.HasNext() {
		currentWord := wordsIter.GetNext()
		maxWordLength = max(maxWordLength, len([]rune(currentWord)))
		cannonicalRepresentation := wordscatalog.SortString(currentWord)
		if equivalentWords, keyExists := equivalenceClasses[cannonicalRepresentation]; keyExists {
			equivalentWords = append(equivalentWords, currentWord)
			equivalenceClasses[cannonicalRepresentation] = equivalentWords
		} else {
			equivalenceClasses[cannonicalRepresentation] = []string{currentWord}
		}
	}
	if err := wordsIter.GetErr(); err != nil {
		return MapBackedWordsCatalog{}, err
	}

	for key, val := range equivalenceClasses {
		// reclaiming any extra memory allocated for the slice, as the length will not change for the duration of the program
		newSlice := make([]string, len(val))
		copy(newSlice, val)
		equivalenceClasses[key] = newSlice
	}

	return MapBackedWordsCatalog{maxWordLength, equivalenceClasses}, nil
}

func (catalog MapBackedWordsCatalog) CountWordsInCatalog() int {
	counter := 0
	for _, equivalentWordsSlice := range catalog.EquivalenceClasses {
		counter += len(equivalentWordsSlice)
	}
	return counter
}
