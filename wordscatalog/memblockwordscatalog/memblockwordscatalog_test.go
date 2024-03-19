package memblockwordscatalog

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

func TestFindEquivalentWords(t *testing.T) {
	catalog := wordscatalog.GetCatalog(t, NewMemBlockWordsCatalog, WORDS_CATALOG_FILE)
	wordscatalog.FindEquivalentWordsTestHelper(t, catalog)
}

func TestNewWordsCatalog(t *testing.T) {
	wordscatalog.NewWordsCatalogTestHelper(t, NewMemBlockWordsCatalog)
}

func TestFindNonExistentWord(t *testing.T) {
	wordscatalog.FindNonExistentWordTestHelper(t, NewMemBlockWordsCatalog)
}

func TestFindWord(t *testing.T) {
	wordscatalog.FindWordTestHelper(t, NewMemBlockWordsCatalog)
}

func TestReadWordsCatalogFromFile(t *testing.T) {
	catalog := wordscatalog.GetCatalog(t, NewMemBlockWordsCatalog, WORDS_CATALOG_FILE)
	wordscatalog.ReadWordsCatalogFromFileTestHelper(t, catalog)
}
