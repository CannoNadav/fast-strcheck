package programargs

import (
	"errors"
	"flag"
	"fmt"
	"main/wordscatalog"
	"main/wordscatalog/mapwordscatalog"
	"main/wordscatalog/memblockwordscatalog"
	"strings"
)

// code in this package should logically have been placed inside the main package,
// but then it wouldn't have been reusable from any other place

const MEMBLOCK_WORDS_CATALOG = "MemBlockWordsCatalog"
const MAP_WORDS_CATALOG = "MapWordsCatalog"

var wordsCatalogAlgorithm *string

func getAlgorithms() map[string]wordscatalog.WordsCatalogFactory {
	return map[string]wordscatalog.WordsCatalogFactory{
		MEMBLOCK_WORDS_CATALOG: memblockwordscatalog.NewMemBlockWordsCatalog,
		MAP_WORDS_CATALOG:      mapwordscatalog.NewMapBackedWordsCatalog,
	}
}

func keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func init() {
	// called on package initialization
	// importing this package will run its init functions, and will parse user args
	wordsCatalogAlgorithm = flag.String("WordsCatalogAlg", "", "The data structure "+
		"to use as to store the similar words catalog. options are: "+strings.Join(keys(getAlgorithms()), ", "))
}

func getDefaultWordsCatalogFactoryName() string {
	//return mapwordscatalog.NewMapBackedWordsCatalog
	return MEMBLOCK_WORDS_CATALOG
}

func GetWordsCatalogFactory(loggerMock func(string)) (string, wordscatalog.WordsCatalogFactory, error) {
	var selectedAlg string
	defaultAlgName := getDefaultWordsCatalogFactoryName()
	if *wordsCatalogAlgorithm == "" {
		loggerMock(fmt.Sprintf("Did not specify datastructure algorithm, using default value: %v", defaultAlgName))
		return defaultAlgName, getAlgorithms()[defaultAlgName], nil
	} else {
		selectedAlg = *wordsCatalogAlgorithm
	}
	for algName, alg := range getAlgorithms() {
		if selectedAlg == algName {
			loggerMock(fmt.Sprintf("Using words catalog algorithm: %v", algName))
			return algName, alg, nil
		}
	}

	return "", nil, errors.New(fmt.Sprintf("Unrecognized value for algorithm: %v. Valid options: %v", *wordsCatalogAlgorithm, strings.Join(keys(getAlgorithms()), ", ")))
}
