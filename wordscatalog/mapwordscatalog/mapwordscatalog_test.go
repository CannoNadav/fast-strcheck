package mapwordscatalog

import (
	"main/wordscatalog"
	"os"
	"testing"
)

const (
	WORDS_CATALOG_PATH = "../../words_clean.txt"
)

var (
	WORDS_CATALOG_FILE *os.File
)

func TestMain(m *testing.M) {
	wordscatalog.RunTestHelper(m, &WORDS_CATALOG_FILE, WORDS_CATALOG_PATH)
}

//func TestPrintRandomEntries(t *testing.T) {
//	catalog := getCatalog(t)
//
//	// Store the map keys in a slice
//	keys := make([]string, 0, len(catalog.EquivalenceClasses))
//	for k := range catalog.EquivalenceClasses {
//		keys = append(keys, k)
//	}
//
//	// Seed the random number generator
//	rand.Seed(time.Now().UnixNano())
//
//	// Select and print random entries from the map
//	numRandomEntries := 15
//	for i := 0; i < numRandomEntries; i++ {
//		randomKey := keys[rand.Intn(len(keys))]
//		fmt.Printf(fmt.Sprintf("Random entry: %v -> %v\n", randomKey, catalog.EquivalenceClasses[randomKey]))
//	}
//}

func TestNewWordsCatalog(t *testing.T) {
	wordscatalog.NewWordsCatalogTestHelper(t, NewMapBackedWordsCatalog)
}

func TestFindNonExistentWord(t *testing.T) {
	wordscatalog.FindNonExistentWordTestHelper(t, NewMapBackedWordsCatalog)
}

func TestFindWord(t *testing.T) {
	wordscatalog.FindWordTestHelper(t, NewMapBackedWordsCatalog)
}

func TestFindEquivalentWords(t *testing.T) {
	catalog := wordscatalog.GetCatalog(t, NewMapBackedWordsCatalog, WORDS_CATALOG_FILE)
	wordscatalog.FindEquivalentWordsTestHelper(t, catalog)
}

func TestReadWordsCatalogFromFile(t *testing.T) {
	catalog := wordscatalog.GetCatalog(t, NewMapBackedWordsCatalog, WORDS_CATALOG_FILE)
	wordscatalog.ReadWordsCatalogFromFileTestHelper(t, catalog)
}
