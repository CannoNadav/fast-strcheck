package optimization

import (
	"fmt"
	"main/wordscatalog"
	"main/wordscatalog/mapwordscatalog"
	"main/wordscatalog/memblockwordscatalog"
	"os"
	"testing"
)

const WORDS_CATALOG_PATH = "../words_clean.txt"

var WORDS_CATALOG_FILE *os.File

//words longer than 0: 351075
//words longer than 1: 351049
//words longer than 2: 350630
//words longer than 3: 348599
//words longer than 4: 341871
//words longer than 5: 327024
//words longer than 6: 298897
//words longer than 7: 259203
//words longer than 8: 210199
//words longer than 9: 159544
//words longer than 10: 116124
//words longer than 11: 80512
//words longer than 12: 52843
//words longer than 13: 32872
//words longer than 14: 19300
//words longer than 15: 10776
//words longer than 16: 5744
//words longer than 17: 2847
//words longer than 18: 1402
//words longer than 19: 653
//words longer than 20: 298
//words longer than 21: 131
//words longer than 22: 58
//words longer than 23: 28
//words longer than 24: 16
//words longer than 25: 8
//words longer than 26: 8
//words longer than 27: 5
//words longer than 28: 3
//words longer than 29: 1
//words longer than 30: 1
//words longer than 31: 0

//there are 0 words of length 0
//there are 26 words of length 1
//there are 419 words of length 2
//there are 2031 words of length 3
//there are 6728 words of length 4
//there are 14847 words of length 5
//there are 28127 words of length 6
//there are 39694 words of length 7
//there are 49004 words of length 8
//there are 50655 words of length 9
//there are 43420 words of length 10
//there are 35612 words of length 11
//there are 27669 words of length 12
//there are 19971 words of length 13
//there are 13572 words of length 14
//there are 8524 words of length 15
//there are 5032 words of length 16
//there are 2897 words of length 17
//there are 1445 words of length 18
//there are 749 words of length 19
//there are 355 words of length 20
//there are 167 words of length 21
//there are 73 words of length 22
//there are 30 words of length 23
//there are 12 words of length 24
//there are 8 words of length 25
//there are 0 words of length 26
//there are 3 words of length 27
//there are 2 words of length 28
//there are 2 words of length 29
//there are 0 words of length 30
//there are 1 words of length 31
//

//Length | Count | % of Total |                             Possible Strings | % of Possible
//------------------------------------------------------------------------------------------------------
//1      |    26 | 0.01%      |                                           26 | 100.000000%
//2 	 |   419 | 0.12% 	  |                                          676 | 61.982249%
//3 	 |  2031 | 0.58% 	  |                                        17576 | 11.555530%
//4 	 |  6728 | 1.92% 	  |                                       456976 | 1.472287%
//5 	 | 14847 | 4.23% 	  |                                     11881376 | 0.124960%
//6 	 | 28127 | 8.01% 	  |                                    308915776 | 0.009105%
//7 	 | 39694 | 11.31% 	  |                                   8031810176 | 0.000494%
//8 	 | 49004 | 13.96% 	  |                                 208827064576 | 0.000023%
//9 	 | 50655 | 14.43% 	  |                                5429503678976 | 0.000001%
//10 	 | 43420 | 12.37% 	  |                              141167095653376 | 0.000000%
//11 	 | 35612 | 10.14% 	  |                             3670344486987776 | 0.000000%
//12 	 | 27669 | 7.88% 	  |                            95428956661682176 | 0.000000%
//13 	 | 19971 | 5.69% 	  |                          2481152873203736576 | 0.000000%
//14  	 | 13572 | 3.87% 	  |                         64509974703297150976 | 0.000000%
//15 	 |  8524 | 2.43% 	  |                       1677259342285725925376 | 0.000000%
//16 	 |  5032 | 1.43% 	  |                      43608742899428874059776 | 0.000000%
//17 	 |  2897 | 0.83% 	  |                    1133827315385150725554176 | 0.000000%
//18 	 |  1445 | 0.41% 	  |                   29479510200013918864408576 | 0.000000%
//19 	 |   749 | 0.21% 	  |                  766467265200361890474622976 | 0.000000%
//20 	 |   355 | 0.10% 	  |                19928148895209409152340197376 | 0.000000%
//21 	 |   167 | 0.05% 	  |               518131871275444637960845131776 | 0.000000%
//22 	 |    73 | 0.02% 	  |             13471428653161560586981973426176 | 0.000000%
//23 	 |    30 | 0.01% 	  |            350257144982200575261531309080576 | 0.000000%
//24 	 |    12 | 0.00% 	  |           9106685769537214956799814036094976 | 0.000000%
//25 	 |     8 | 0.00% 	  |         236773830007967588876795164938469376 | 0.000000%
//26     |     0 | 0.00%      |        6156119580207157310796674288400203776 | 0.000000%
//27 	 |     3 | 0.00% 	  |      160059109085386090080713531498405298176 | 0.000000%
//28 	 |     2 | 0.00% 	  |     4161536836220038342098551818958537752576 | 0.000000%
//29 	 |     2 | 0.00% 	  |   108199957741720996894562347292921981566976 | 0.000000%
//30     |     0 | 0.00%      |  2813198901284745919258621029615971520741376 | 0.000000%
//31 	 |     1 | 0.00% 	  | 73143171433403393900724146770015259539275776 | 0.000000%
//
//
// Amount of Equivalence classes by length:
// conclusion: almost all words have 0 equivalent words (are in equivalence class 1)

//There are 283543 equivalence classes with length 1
//There are 20847 equivalence classes with length 2
//There are 1486 equivalence classes with length 4
//There are 4630 equivalence classes with length 3
//There are 31 equivalence classes with length 9
//There are 558 equivalence classes with length 5
//There are 132 equivalence classes with length 7
//There are 237 equivalence classes with length 6
//There are 5 equivalence classes with length 11
//There are 1 equivalence classes with length 13
//There are 45 equivalence classes with length 8
//There are 8 equivalence classes with length 10
//There are 1 equivalence classes with length 15
//There are 3 equivalence classes with length 14
//There are 2 equivalence classes with length 12

// The equivalence classes of length 1 (each just a word) are characterized by:
//There are 26 words of length 1
//There are 115 words of length 2
//There are 650 words of length 3
//There are 2489 words of length 4
//There are 6781 words of length 5
//There are 15252 words of length 6
//There are 26023 words of length 7
//There are 37423 words of length 8
//There are 42720 words of length 9
//There are 39573 words of length 10
//There are 33817 words of length 11
//There are 26719 words of length 12
//There are 19596 words of length 13
//There are 13373 words of length 14
//There are 8400 words of length 15
//There are 4954 words of length 16
//There are 2847 words of length 17
//There are 1419 words of length 18
//There are 735 words of length 19
//There are 349 words of length 20
//There are 157 words of length 21
//There are 67 words of length 22
//There are 30 words of length 23
//There are 12 words of length 24
//There are 8 words of length 25
//There are 3 words of length 27
//There are 2 words of length 28
//There are 2 words of length 29
//There are 1 words of length 31

// if we replace each word in the catalog with its hash each word will take 8 bytes (one machine word)

//>>> words_length_in_equiv_cls_one
//{1: 26, 2: 115, 3: 650, 6: 15252, 4: 2489, 5: 6781, 8: 37423, 9: 42720, 10: 39573, 7: 26023, 11: 33817, 12: 26719, 14: 13373, 13: 19596, 16: 4954, 15: 8400, 20: 349, 19: 735, 17: 2847, 18: 1419, 21: 157, 22: 67, 23: 30, 25: 8, 28: 2, 29: 2, 31: 1, 24: 12, 27: 3}
//>>>
//>>>
//>>>
//>>>
//>>> sum([length * count for length, count in words_length_in_equiv_cls_one.items()])
//2840244
//>>>
//>>>
//>>> sum([8*count for count in words_length_in_equiv_cls_one.values()]
//... )
//2268344
//>>>
//>>>
//>>>
//>>> current_memory_size = sum([length * count for length, count in words_length_in_equiv_cls_one.items()])
//>>>
//>>>
//>>> optimized_memory_size = sum([8*count for count in words_length_in_equiv_cls_one.values()])
//>>>
//>>>
//>>>
//>>> optimized_memory_size / current_memory_size
//0.7986440601582118
//
// ~20% memory savings is not worth the added complexity..
// if the average string would have been longer that might have been a useful optimization.

// going for a different approach, for every 1<= i <= 31 all strings of length i are stored in one contiguous
// block of memory. since all letters in the words catalog are ascii lowercase they take just one byte each
// (all ascii characters take one byte in UTF-8).
// Would have used an array but arrays in to be compile time declared with their size, so will use a
// byte block for that. all Strings in the "array" are sorted with the comparison key being their cannonical representation.
// thus all equivalent keys will be stored right next to each other, and all there is to do to pick them up is keep going one by
// one until you find one that is not equivalent to the query string.

// however since there are quite a lot of keys doing a binary search every time might take a few jumps to find..
//>>> math.log(351075, 2)
//18.421419740207835

// so I will have a separate array for each string length and each first letter. that way some jumps would be saved

//max bucket for length 1: ('a', 1). log(1, 2)=0.0
//max bucket for length 2: ('a', 23). log(23, 2)=4.523561956057013
//max bucket for length 3: ('a', 161). log(161, 2)=7.330916878114618
//max bucket for length 4: ('s', 597). log(597, 2)=9.221587121264806
//max bucket for length 5: ('s', 1716). log(1716, 2)=10.744833837499547
//max bucket for length 6: ('s', 3040). log(3040, 2)=11.569855608330947
//max bucket for length 7: ('s', 4345). log(4345, 2)=12.085140461701764
//max bucket for length 8: ('s', 5695). log(5695, 2)=12.475480126595476
//max bucket for length 9: ('s', 5616). log(5616, 2)=12.45532722030456
//max bucket for length 10: ('s', 4512). log(4512, 2)=12.139551352398795
//max bucket for length 11: ('p', 4003). log(4003, 2)=11.966865900387539
//max bucket for length 12: ('p', 3218). log(3218, 2)=11.651948610723448
//max bucket for length 13: ('p', 2562). log(2562, 2)=11.323054760341646
//max bucket for length 14: ('p', 1725). log(1725, 2)=10.752380646552893
//max bucket for length 15: ('p', 1116). log(1116, 2)=10.124121311829187
//max bucket for length 16: ('p', 674). log(674, 2)=9.396604781181859
//max bucket for length 17: ('p', 393). log(393, 2)=8.618385502258606
//max bucket for length 18: ('p', 209). log(209, 2)=7.7073591320808825
//max bucket for length 19: ('p', 112). log(112, 2)=6.807354922057604
//max bucket for length 20: ('p', 70). log(70, 2)=6.129283016944967
//max bucket for length 21: ('p', 28). log(28, 2)=4.807354922057604
//max bucket for length 22: ('p', 12). log(12, 2)=3.5849625007211565
//max bucket for length 23: ('p', 9). log(9, 2)=3.1699250014423126
//max bucket for length 24: ('p', 3). log(3, 2)=1.5849625007211563
//max bucket for length 25: ('a', 1). log(1, 2)=0.0
//max bucket for length 27: ('e', 1). log(1, 2)=0.0
//max bucket for length 28: ('a', 1). log(1, 2)=0.0
//max bucket for length 29: ('c', 1). log(1, 2)=0.0
//max bucket for length 31: ('d', 1). log(1, 2)=0.0

//func stringMemoryUsage(s string) int {
//	return len(s) + int(unsafe.Sizeof(s))
//}
//
//func sliceOfStringMemoryUsage(slice []string) int {
//	// Memory taken by the slice header itself
//	memUsage := int(unsafe.Sizeof(slice))
//
//	// Memory used by the strings in the slice
//	for _, s := range slice {
//		memUsage += stringMemoryUsage(s)
//	}
//
//	return memUsage
//}

//func mapMemoryUsage(m map[string][]string) int {
//	totalMemory := 0
//	totalMemory += unsafe.Sizeof(&m)
//
//	for k, v := range m {
//		totalMemory += stringMemoryUsage(k)        // Memory for the key
//		totalMemory += sliceOfStringMemoryUsage(v) // Memory for the value (which is a slice of strings)
//	}
//
//	return totalMemory
//}
//
//func setMemoryUsage(m map[string][]string) int {
//	totalMemory := 0
//
//	for k, v := range m {
//		totalMemory += stringMemoryUsage(k)        // Memory for the key
//		totalMemory += sliceOfStringMemoryUsage(v) // Memory for the value (which is a slice of strings)
//	}
//
//	return totalMemory
//}

func TestMain(m *testing.M) {
	var err error
	WORDS_CATALOG_FILE, err = os.Open(WORDS_CATALOG_PATH)
	if err != nil || WORDS_CATALOG_FILE == nil {
		panic(fmt.Sprintf("Test Initialization failed: Cannot open the words catalog file '%s', %s ", WORDS_CATALOG_PATH, err.Error()))
	}
	code := m.Run() // This will run the actual tests
	defer func() {
		if closeErr := WORDS_CATALOG_FILE.Close(); closeErr != nil {
			fmt.Printf("Failed to close the words catalog file: %v", closeErr)
		}
		os.Exit(code)
	}()
}

func deepCopyString(str string) string {
	byteSlice := []byte(str)
	copiedBytes := make([]byte, len(str))
	copy(copiedBytes, byteSlice)
	return string(copiedBytes)
}

func findMapMemorySize(targetMap map[string][]string, t *testing.T, benchmarkName string) uint64 {
	// trying to force allocations on the heap to be able to count them
	var res testing.BenchmarkResult
	var copiedSlice []string
	var copiedKey string
	var byteSlice []byte
	var copiedBytes []byte
	var copiedString string
	var copiedMapReferenceOutsideItsDeclaringClosure map[string][]string // helps the copied map outlive its function's scope

	t.Run(benchmarkName, func(t *testing.T) {
		res = testing.Benchmark(func(b *testing.B) {
			b.ReportAllocs()
			// putting rest of code into a closure to verify that no interesting(tested) code will not be compile time reordered to happen before ReportAllocs is called
			// the map is returned from the closure to verify the variable will be allocated on the heap (needs to outlive the closure's scope)
			copiedMapReferenceOutsideItsDeclaringClosure =
				func() map[string][]string {
					var copiedMap map[string][]string

					copiedMap = make(map[string][]string)
					for key, originalSlice := range targetMap {
						if originalSlice == nil {
							copiedSlice = nil
						} else {
							copiedSlice = make([]string, cap(originalSlice))
							for _, str := range originalSlice {
								copiedString = deepCopyString(str)
								copiedSlice = append(copiedSlice, copiedString)
							}
						}

						byteSlice = []byte(key)
						copiedBytes = make([]byte, len(key))
						copy(copiedBytes, byteSlice)
						copiedKey = string(copiedBytes)
						copiedMap[copiedKey] = copiedSlice
					}

					return copiedMap
				}()
		})
	})

	if copiedMapReferenceOutsideItsDeclaringClosure == nil {
		panic("The target map was not copied, or an internal error occurred")
	}
	return res.MemBytes
}

//func findSetMemorySize(targetSet datastructures.Set, t *testing.T, benchmarkName string) uint64 {
//	// to find the memory size of a set (with the given implementation) I make its underlying list a valid
//	// target for findMapMemorySize and then substructs the amount of additional memory added
//
//}

func Test_CatalogSizeInDifferentImplementation(t *testing.T) {

	if WORDS_CATALOG_FILE == nil {
		t.Fatalf("Words catalog1 file is not initialized")
	}
	_, err := WORDS_CATALOG_FILE.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to reset file pointer: %equivClass", err)
	}
	catalog1, err := wordscatalog.ReadWordsCatalogFromFile(WORDS_CATALOG_FILE, mapwordscatalog.NewMapBackedWordsCatalog)
	if err != nil {
		t.Fatalf("failed to create the words catalog1 after file was open for read: " + err.Error())
	}
	mapWordsCatalog := catalog1.(mapwordscatalog.MapBackedWordsCatalog)

	totalBytes := findMapMemorySize(mapWordsCatalog.EquivalenceClasses, t, "Original mapWordsCatalog with all words")
	fmt.Printf("Total Memory used by the mapWordsCatalog: %d bytes (%.2f MB)\n", totalBytes, float64(totalBytes)/1024/1024)

	mapWithNullSlices := make(map[string][]string)

	for k, v := range mapWordsCatalog.EquivalenceClasses {
		if len(v) == 1 {
			mapWithNullSlices[k] = nil
		} else {
			mapWithNullSlices[k] = v
		}
	}

	totalBytes = findMapMemorySize(mapWithNullSlices, t, "mapWordsCatalog with all 1 sized word slices nil'd")
	fmt.Printf("Total Memory used by the mapWordsCatalog, after shrinking: %d bytes (%.2f MB)\n", totalBytes, float64(totalBytes)/1024/1024)

	mapWithout1SizedValues := make(map[string][]string)
	//setOfOneSizedValues := datastructures.NewSet()

	memSizeFor1SizedEquivClassesStrings := 0
	memSizeForRestOfEquivClassesStrings := 0 // this represents total memory in case we keep all strings in the 1-block array
	for key, equivClass := range mapWordsCatalog.EquivalenceClasses {
		if len(equivClass) == 1 {
			memSizeFor1SizedEquivClassesStrings += len([]rune(equivClass[0]))
			//setOfOneSizedValues.Add(equivClass[0])
		} else {
			for _, str := range equivClass {
				memSizeForRestOfEquivClassesStrings += len([]rune(str))
			}
			mapWithout1SizedValues[key] = equivClass
		}
	}

	totalBytesMap := findMapMemorySize(mapWithout1SizedValues, t, "mapWordsCatalog with all 1 sized word slices completely removed")
	//totalBytesSet := findMapMemorySize(setOfOneSizedValues.EmbeddedMap, , t, "mapWordsCatalog with all 1 sized word slices completely removed"))
	fmt.Printf("Total Memory used by the mapWordsCatalog, without 1 sized equiv classes: %d bytes (%.2f MB)\n", totalBytesMap, float64(totalBytesMap)/1024/1024)
	fmt.Printf("Total additional Memory needed to store the 1 sized equiv classes: %d bytes (%.2f MB)\n", memSizeFor1SizedEquivClassesStrings, float64(memSizeFor1SizedEquivClassesStrings)/1024/1024)
	fmt.Printf("Total additional Memory needed if we also store the rest of equiv classes strings in contiguous block: %d bytes (%.2f MB)\n", memSizeForRestOfEquivClassesStrings, float64(memSizeForRestOfEquivClassesStrings)/1024/1024)
	fmt.Printf("Total theoretical Memory needed for all strings in the single mem block approach: %d bytes (%.2f MB)\n", memSizeForRestOfEquivClassesStrings+memSizeFor1SizedEquivClassesStrings, float64(memSizeForRestOfEquivClassesStrings+memSizeFor1SizedEquivClassesStrings)/1024/1024)

	_, err = WORDS_CATALOG_FILE.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to reset file pointer: %equivClass", err)
	}
	catalog2, err := wordscatalog.ReadWordsCatalogFromFile(WORDS_CATALOG_FILE, memblockwordscatalog.NewMemBlockWordsCatalog)
	if err != nil {
		t.Fatalf("failed to create the words catalog1 after file was open for read: " + err.Error())
	}
	memBlockWordsCatalog := catalog2.(memblockwordscatalog.MemBlockWordsCatalog)
	memBlockMemorySize := memBlockWordsCatalog.GetMemorySize()
	fmt.Printf("Actual memory usage for MemBlockWordsCatalog: %d bytes (%.2f MB)\n", memBlockMemorySize, float64(memBlockMemorySize)/1024/1024)
}
