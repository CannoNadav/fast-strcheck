package test_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/server"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const (
	WORDS_CATALOG_FILE_PATH = "../words_clean.txt"
)

type SimilarWordsResponse struct {
	SimilarWords []string `json:"similar"`
}

type StatsResponse struct {
	TotalWords          int `json:"totalWords"`
	TotalRequests       int `json:"totalRequests"`
	AvgProcessingTimeNs int `json:"avgProcessingTimeNs"`
}

func (statsResponse StatsResponse) String() string {
	return fmt.Sprintf("StatsResponse: (TotalRequests=%v, AvgProcessingTimeNs=%v, TotalWords=%v)", statsResponse.TotalRequests, statsResponse.AvgProcessingTimeNs, statsResponse.TotalWords)
}

func sendReq(t *testing.T, url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer func() {
		if resp != nil {
			err2 := resp.Body.Close()
			if err2 != nil {
				t.Log("Couldn't close response for statsResponse request(not failing test for that)")
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Received non-OK HTTP status for get request to url: %v", url))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read response from endpoint: %v. msg: %v", url, err))
	}

	return body, nil
}

func SendStatsReq(t *testing.T, numStatsRequestsFinished *atomic.Int32, numWordsInCatalog int) (StatsResponse, bool) {
	// Make request to statsResponse endpoint and check the response validity
	// notice we can't assert the number of expected similar words requests because of race conditions from network etc

	var statsResponse StatsResponse
	body, err := sendReq(t, "http://localhost"+server.SERVER_ADDRESS+server.STATS_ENDPOINT_PATH)
	if err != nil {
		t.Error(err)
		return statsResponse, true
	}

	if err = json.Unmarshal(body, &statsResponse); err != nil {
		t.Error(fmt.Sprintf("Failed to unmarshal stats response: %v", err))
		return statsResponse, true
	}

	numStatsRequestsFinished.Add(1)

	if statsResponse.TotalWords != numWordsInCatalog {
		t.Error("received wrong number of words in catalog")
		return statsResponse, true
	}

	if statsResponse.AvgProcessingTimeNs == 0 && statsResponse.TotalRequests != 0 || statsResponse.AvgProcessingTimeNs != 0 && statsResponse.TotalRequests == 0 {
		t.Error(fmt.Sprintf("avg processing is zero iff number of requests is zero.  AvgProcessingTimeNs: %v, total similar words requests: %v", statsResponse.AvgProcessingTimeNs, statsResponse.TotalRequests))
		return statsResponse, true
	}

	return statsResponse, false
}

func IsServerReady(t *testing.T) bool {
	url := "http://localhost" + server.SERVER_ADDRESS + "/prob-server-ready"
	body, err := sendReq(t, url)
	return err == nil && string(body) == "ok"
}

func getRandomStringFromSlice(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	return strs[rand.Intn(len(strs))]
}

func __sendSimilarWordsRequest(t *testing.T, allWordsInCatalog []string, numSimilarWordsRequestsFinished *atomic.Int32) (bool, []byte) {
	searchWord := getRandomStringFromSlice(allWordsInCatalog)
	body, err := sendReq(t, "http://localhost"+server.SERVER_ADDRESS+server.SIMILAR_WORDS_ENDPOINT_PATH+"?word="+searchWord)

	if err != nil {
		t.Error(fmt.Sprintf("Failed to make similar words request with word %v: %v", searchWord, err))
		return false, nil
	}

	numSimilarWordsRequestsFinished.Add(1)
	return true, body
}

func sendSimilarWordsRequestAndUnmarshall(t *testing.T, word string) (SimilarWordsResponse, bool) {
	var numSimilarWordsRequestsFinished atomic.Int32
	success, body := __sendSimilarWordsRequest(t, []string{word}, &numSimilarWordsRequestsFinished)
	var response SimilarWordsResponse
	if !success {
		return response, false
	}
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Error("could not unmarshall request body")
		return response, false
	}
	return response, true
}

func sendSimilarWordsRequestWithoutUnmarshalling(t *testing.T, allWordsInCatalog []string, numSimilarWordsRequestsFinished *atomic.Int32) bool {
	success, _ := __sendSimilarWordsRequest(t, allWordsInCatalog, numSimilarWordsRequestsFinished)
	return success
}

func waitForServerReady(t *testing.T) {
	// Periodically check for server readiness
	t.Log("Waiting for server to become available, to start the test")

	for ; !IsServerReady(t); time.Sleep(50 * time.Millisecond) { // Wait for 50ms before retrying
	}
}

func getAllWordsInCatalog(t *testing.T) []string {
	// reads the entire file into memory.. since it's just a test it's bearable
	data, err := os.ReadFile(WORDS_CATALOG_FILE_PATH)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to read file: %s", err))
	}

	return strings.Split(string(data), "\n")
}

func sendShutdownSignalToServer(t *testing.T, didTestFail bool) {
	var result string
	if didTestFail {
		result = "Fail"
	} else {
		result = "Pass"
	}

	t.Log("signaling server to end test")
	resp, err := http.Get("http://localhost" + server.SERVER_ADDRESS + "/stop-clock-and-shutdown?testResult=" + result)
	if err != nil {
		t.Fatalf("Did not manage to send the shutdown signal to the server: " + err.Error())
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("Error closing response body: " + err.Error())
	}
}

func sendStartSignalToServer(t *testing.T) {
	//client := &http.Client{
	//	Timeout: 10 * time.Second,
	//}
	resp, err := http.Get("http://localhost" + server.SERVER_ADDRESS + "/start-clock")
	if err != nil {
		t.Fatalf("Did not manage to send the start signal to the server: " + err.Error())
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("Error closing response body: " + err.Error())
	}
}

func RunTest(t *testing.T, testName string, testTarget func(*testing.T, []string) bool) {
	allWordsInCatalog := getAllWordsInCatalog(t)
	waitForServerReady(t)
	t.Log(fmt.Sprintf("Server is ready, starting test: %v", testName))
	sendStartSignalToServer(t)
	testStartTime := time.Now()
	didTestFail := testTarget(t, allWordsInCatalog)
	testEndTime := time.Now()
	sendShutdownSignalToServer(t, didTestFail)
	t.Log(fmt.Sprintf("Client: Test finished in %v seconds\n", testEndTime.Sub(testStartTime).Seconds()))
	if didTestFail {
		t.Fatalf("Client: Test failed, see logs for details")
	}
}
