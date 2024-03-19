package wordscatalog

import (
	"fmt"
	"log"
	"main/datastructures"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	expectedNumOfWordsInCatalog = 351075
	longestWordInCatalog        = "dichlorodiphenyltrichloroethane"
	expectedNumOfEquivClasses   = 311529
)

func RunTestHelper(m *testing.M, wordsCatalogFile **os.File, WordsCatalogFilePath string) {
	var err error
	*wordsCatalogFile, err = os.Open(WordsCatalogFilePath)
	if err != nil || wordsCatalogFile == nil {
		panic(fmt.Sprintf("Test Initialization failed: Cannot open the words catalog file '%s', %s ", WordsCatalogFilePath, err.Error()))
	}
	code := m.Run() // This will run the actual tests
	defer func() {
		if closeErr := (*wordsCatalogFile).Close(); closeErr != nil {
			log.Fatalf("Failed to close the words catalog file: %v", closeErr)
		}
		os.Exit(code)
	}()
}

func GetCatalog(t *testing.T, factory WordsCatalogFactory, wordsCatalogFile *os.File) WordsCatalog {
	if wordsCatalogFile == nil {
		t.Fatalf("Words catalog file is not initialized")
	}

	_, err := wordsCatalogFile.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to reset file pointer: %v", err)
	}
	catalog, err := ReadWordsCatalogFromFile(wordsCatalogFile, factory)
	if err != nil {
		t.Fatalf("failed to create the words catalog after file was open for read: " + err.Error())
	}
	return catalog
}

func FindEquivalentWordsTestHelper(t *testing.T, catalog WordsCatalog) {

	tests := []struct {
		inputWord string
		expected  []string
	}{
		{"elaterins", []string{"entailers", "intersale", "larsenite", "nearliest", "treenails"}},
		{"morticians", []string{"romanistic"}},
		{"aspirants", []string{"partisans"}},
		{"myeloproliferative", []string{}}, // empty because no other equivalent words
		{"nu", []string{"un"}},
		{"monroe", []string{"mooner", "morone"}},
		{"arow", []string{}}, // empty because no other equivalent words
		{"moonsets", []string{"mootness"}},
		{"dorsers", []string{"drosser"}},
		{"eker", []string{"erke", "reek"}},
		//{"pstarisan", []string{"aspirants", "partisans"}},
		{"nspaisart", []string{"aspirants", "partisans"}},
	}

	for _, test := range tests {
		t.Run(test.inputWord, func(t *testing.T) {

			fmt.Printf("Searching for word: %s\n", test.inputWord)
			result := catalog.FindEquivalentWords(test.inputWord)

			// Sort the slices for easier comparison
			sort.Strings(result)
			sort.Strings(test.expected)

			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("For word %s, expected %v, got %v", test.inputWord, test.expected, result)
			}
			println(fmt.Sprintf("OK!: For word %s, expected %v, got %v\n", test.inputWord, test.expected, result))
		})
	}
}

func NewWordsCatalogTestHelper(t *testing.T, factory WordsCatalogFactory) {
	scanner := datastructures.NewSliceIterator([]string{"apple", "papel", "dog", "god"})
	catalog, err := factory(scanner)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	numEquivalenceClasses := catalog.GetNumEquivalenceClasses()
	if numEquivalenceClasses != 2 {
		t.Fatalf("Expected 2 equivalence classes, got %d", numEquivalenceClasses)
	}

	result := catalog.FindEquivalentWords("apple")
	if len(result) != 1 {
		t.Fatalf("Expected 1 equivalent words for apple, got %d", len(result))
	}
}

func FindNonExistentWordTestHelper(t *testing.T, factory WordsCatalogFactory) {
	scanner := datastructures.NewSliceIterator([]string{"apple", "papel"})
	catalog, err := factory(scanner)
	if err != nil {
		t.Fatalf("failed creating a words catalog: " + err.Error())
	}
	match := catalog.FindEquivalentWords("chinchilla")
	if len(match) != 0 {
		t.Fatalf("Expected empty match, got %v", err)
	}
}

func FindWordTestHelper(t *testing.T, factory WordsCatalogFactory) {
	tests := []struct {
		name     string
		words    []string
		search   string
		expected []string
	}{
		{
			name:     "word without spaces in anagrams",
			words:    []string{"enthusiastic", "chastenitusi", "entuhitssica", "enusthastici"},
			search:   "enthusiastic",
			expected: []string{"chastenitusi", "entuhitssica", "enusthastici"},
		},
		{
			name:     "very long nonexistent word",
			words:    []string{"accommodate", "a mood came t", "dalmatian eco"},
			search:   "pneumonoultramicroscopicsilicovolcanoconiosis",
			expected: []string{},
		},
		{
			name:     "empty word",
			words:    []string{"resplendent", "spend pert l", "dents repel"},
			search:   "",
			expected: []string{},
		},
		{
			name:     "anagrams with spaces",
			words:    []string{"listen", "silent", "en list", "lens it", "panacea"},
			search:   "en list",
			expected: []string{"lens it"},
		},
		{
			name:     "word longer than max allowed",
			words:    []string{"short", "torch", "orts h"},
			search:   "exceptionallylongwordthatexceedsmaxlimit",
			expected: []string{},
		},
		{
			name:     "word with special characters",
			words:    []string{"$dog", "god$", "g$od"},
			search:   "g$od",
			expected: []string{"$dog", "god$"},
		},
		{
			name:     "word with numbers",
			words:    []string{"go1d", "1god", "d1og"},
			search:   "1god",
			expected: []string{"go1d", "d1og"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := datastructures.NewSliceIterator(tt.words)
			catalog, err := factory(scanner)
			if err != nil {
				t.Fatalf("Unexpected error when setting up words catalog: %v", err)
			}

			result := catalog.FindEquivalentWords(tt.search)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}

func shuffleString(s string) string {
	for {
		rand.Seed(time.Now().UnixNano())
		runeArray := []rune(s)
		n := len(runeArray)
		for i := n - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			runeArray[i], runeArray[j] = runeArray[j], runeArray[i]
		}
		res := string(runeArray)
		if res != s {
			return string(runeArray)
		}
	}
}

func ReadWordsCatalogFromFileTestHelper(t *testing.T, catalogFromFile WordsCatalog) {
	numWordsInCatalog := catalogFromFile.CountWordsInCatalog()
	if numWordsInCatalog != expectedNumOfWordsInCatalog {
		t.Fatalf("unexpected number of words in catalog. expected: %d, found: %d\n", expectedNumOfWordsInCatalog, numWordsInCatalog)
	}

	numEquivalenceClasses := catalogFromFile.GetNumEquivalenceClasses()
	if numEquivalenceClasses != expectedNumOfEquivClasses {
		t.Fatalf("unexpected number of equivalence classes. expected: %d, found: %d\n", expectedNumOfEquivClasses, numEquivalenceClasses)
	}

	//maxWordLength := 0
	//for key, equivalentWordsSlice := range catalog.EquivalenceClasses {
	//	for _, word := range equivalentWordsSlice {
	//		if chars := []rune(word); len(chars) > maxWordLength {
	//			maxWordLength = len(chars)
	//		}
	//		if SortString(word) != key {
	//			t.Fatalf("word %s was in bucket %s, but the two contain different characters", word, key)
	//		}
	//	}
	//}

	resForLongestWord := catalogFromFile.FindEquivalentWords(longestWordInCatalog)
	resForShuffeledLongestWord := catalogFromFile.FindEquivalentWords(shuffleString(longestWordInCatalog))
	if len(resForShuffeledLongestWord) != 1 {
		t.Fatalf("expected one word in equiv group, got %v", len(resForShuffeledLongestWord))
	}
	if len(resForLongestWord) != 0 {
		t.Fatalf("expected zero words in equiv group, got %v", len(resForShuffeledLongestWord))
	}

	//catalogFromFile.FindEquivalentWords(longestWordInCatalog)
	//expectedLength := len([]rune(longestWordInCatalog))
	//if maxWordLength != expectedLength {
	//	t.Fatalf("unexpected result for longest word length. expected: %v, actual: %v", expectedLength, maxWordLength)
	//}
}
