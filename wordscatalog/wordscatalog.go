package wordscatalog

import (
	"main/datastructures"
	"os"
	"sort"
)

type WordsCatalogFactory = func(iterator datastructures.Iterator[string]) (WordsCatalog, error)

type WordsCatalog interface {
	CountWordsInCatalog() int
	FindEquivalentWords(word string) []string
	GetNumEquivalenceClasses() int
}

func SortString(s string) string {
	chars := []rune(s)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	return string(chars)
}

func ReadWordsCatalogFromFile(catalogFile *os.File, catalogFactory WordsCatalogFactory) (WordsCatalog, error) {
	lineIterator := datastructures.NewFileLineIterator(catalogFile)
	catalog, err := catalogFactory(lineIterator)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}
