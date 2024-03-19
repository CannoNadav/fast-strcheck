package test_client

import (
	"fmt"
	"io"
	"main/server"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func SortString(s string) string {
	chars := []rune(s)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	return string(chars)
}

func findEquivalenceClasses(allWordsInCatalog []string) map[string][]string {
	equivalenceClasses := make(map[string][]string)
	for _, currentWord := range allWordsInCatalog {
		cannonicalRepresentation := SortString(currentWord)
		if equivalentWords, keyExists := equivalenceClasses[cannonicalRepresentation]; keyExists {
			equivalentWords = append(equivalentWords, currentWord)
			equivalenceClasses[cannonicalRepresentation] = equivalentWords
		} else {
			equivalenceClasses[cannonicalRepresentation] = []string{currentWord}
		}
	}

	return equivalenceClasses
}

func generateRandomASCIIString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomWordsTestTarget(t *testing.T, allWordsInCatalog []string) bool {
	const numGoroutines = 100
	const numRequestsPerGoroutine = 1000
	const maxStrLength = 16
	var wg sync.WaitGroup
	var didTestFail atomic.Bool // completely unnecessary to be atomic atm since there are only concurrent writes are of the same value(true) but better still practice..
	didTestFail.Store(false)

	equivClasses := findEquivalenceClasses(allWordsInCatalog)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			var searchWord string
			for j := 0; j < numRequestsPerGoroutine; j++ {
				index := rand.Intn(int(float64(len(allWordsInCatalog)) * 1.1)) // gives a 10% change for totally random words (they almost always have no result)
				if index < len(equivClasses) {
					searchWord = allWordsInCatalog[index]
				} else {
					searchWord = generateRandomASCIIString(rand.Intn(maxStrLength) + 1)
				}
				expectedResult := []string(nil)
				equivClass, ok := equivClasses[SortString(searchWord)]
				if ok {
					for _, s := range equivClass {
						if searchWord != s {
							expectedResult = append(expectedResult, s)
						}
					}
				}

				similarWordsResponse, ok := sendSimilarWordsRequestAndUnmarshall(t, searchWord)
				response := similarWordsResponse.SimilarWords
				if ok {
					if !(len(expectedResult) == 0 && len(response) == 0) { // works around the []string{} != []string(nil) issue
						sort.Strings(similarWordsResponse.SimilarWords)
						sort.Strings(expectedResult)
						if !reflect.DeepEqual(similarWordsResponse.SimilarWords, expectedResult) {
							t.Errorf("Unexpected result for word %v expected response: %v, actual response: %v", searchWord, expectedResult, similarWordsResponse.SimilarWords)
							didTestFail.Store(true)
						}
					}
				} else {
					t.Errorf("Request for word: %v failed", searchWord)
					didTestFail.Store(true)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	t.Log(fmt.Sprintf("Finished sending and validating response correctness for %v similar words requests made of both random strings(max legnth = %v), and words from the catalog.", numGoroutines*numRequestsPerGoroutine, maxStrLength))
	return didTestFail.Load()
}

func TestResponseValidatingTestClient(t *testing.T) {
	RunTest(t, "Random words test with deep validation of response", randomWordsTestTarget)
}

func logSimilarWordsResponse(t *testing.T, word string) {
	var numSimilarWordsRequestsFinished atomic.Int32
	success, body := __sendSimilarWordsRequest(t, []string{word}, &numSimilarWordsRequestsFinished)
	if !success {
		t.Error("Similar words req failed. word: " + word)
	} else {
		t.Log("Similar words req succeeded. word: " + word + "\n" + string(body))
	}
}

func logStatsResponse(t *testing.T) {

	resp, err := http.Get("http://localhost" + server.SERVER_ADDRESS + server.STATS_ENDPOINT_PATH)
	defer func() {
		if resp != nil {
			err2 := resp.Body.Close()
			if err2 != nil {
				t.Error("Couldn't close response for statsResponse request(not failing test for that)")
			}
		}
	}()

	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to make request to statsResponse endpoint: %v, type: . error type: %v", err, reflect.TypeOf(err)))
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response from statsResponse endpoint: %v", err)
		return
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Received non-OK HTTP status from stats endpoint: %s", body)
		return
	}
	t.Log("Stats request succeeded:\n" + string(body))
}

func TestSanity(t *testing.T) {
	RunTest(t, "Log three similar words requests and one stats requests, to see the output", func(t *testing.T, strings []string) bool {
		logSimilarWordsResponse(t, "ppale")        // exists in catalog
		logSimilarWordsResponse(t, "enthusiastic") // exists in catalog
		logSimilarWordsResponse(t, "fwepokmkweic") // doesn't exist in catalog

		//resp, _ := http.Get("http://localhost" + server.SERVER_ADDRESS + server.SIMILAR_WORDS_ENDPOINT_PATH + "?wo")
		//body, _ := io.ReadAll(resp.Body)
		//_ = resp.Body.Close()
		//if string(body) != "{\"error\":\"word parameter is missing\"}" {
		//	t.Fatal("bad response for similar words req with missing word param: " + string(body))
		//}

		logStatsResponse(t)

		return false
	})
}
