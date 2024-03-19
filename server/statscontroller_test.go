package server

import (
	"encoding/json"
	"fmt"
	"io"
	"main/programargs"
	"net/http"
	"testing"
	"time"
)

type StatsResponse struct {
	TotalWords          int `json:"totalWords"`
	TotalRequests       int `json:"totalRequests"`
	AvgProcessingTimeNs int `json:"avgProcessingTimeNs"`
}

func sendSimilarWordsReq(t *testing.T) bool {
	resp, err1 := http.Get("http://localhost" + SERVER_ADDRESS + SIMILAR_WORDS_ENDPOINT_PATH + "?word=apple")
	if resp != nil {
		// Close the response body to avoid resource leak
		err2 := resp.Body.Close()
		if err2 != nil {
			t.Error(fmt.Sprintf("failed to close response body, %v", err2))
		}
	}

	return err1 == nil && resp.StatusCode == http.StatusOK
}

func sendStatsReq(t *testing.T) (StatsResponse, bool) {
	// Make request to stats endpoint
	resp, err := http.Get("http://localhost" + SERVER_ADDRESS + STATS_ENDPOINT_PATH)
	if err != nil {
		t.Fatalf("Failed to make request to stats endpoint: %v", err)
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			t.Errorf(fmt.Sprintf("Failed to close response body(not failing test): %v", err))
		}
	}()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response from stats endpoint: %v", err)
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Received non-OK HTTP status from stats endpoint: %s", body)
	}

	var stats StatsResponse
	if err = json.Unmarshal(body, &stats); err != nil {
		t.Fatalf("Failed to unmarshal stats response: %v", err)
	}
	return stats, true
}

func waitForServerReady(t *testing.T) {
	// Periodically check for server readiness
	for {
		resp, err := http.Get("http://localhost" + SERVER_ADDRESS + SIMILAR_WORDS_ENDPOINT_PATH + "?word=apple")
		if resp != nil {
			// Close the response body to avoid resource leak
			err = resp.Body.Close()
			if err != nil {
				t.Error(fmt.Sprintf("failed to close response body, %v", err))
			}
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(200 * time.Millisecond) // Wait for 200ms before retrying
	}
}

func TestStatsController(t *testing.T) {
	statsInfoChannel := make(chan StatsInfo, CHANNEL_SIZE)
	cpuTimeChannel := make(chan int, CHANNEL_SIZE)
	defer func() {
		close(statsInfoChannel)
		close(cpuTimeChannel)
	}()

	_, wordsCatalogFactory, _ := programargs.GetWordsCatalogFactory(func(s string) {
		t.Log(s)
	})
	similarWordsService := InitNewServer(cpuTimeChannel, statsInfoChannel, wordsCatalogFactory, WORDS_CATALOG_FILE_PATH)
	go similarWordsService.StartServer()
	defer similarWordsService.StopServer()

	waitForServerReady(t)
	statsResponse1, ok := sendStatsReq(t)
	if !ok {
		t.Fatalf("Bad response for stats request")
	}
	if statsResponse1.TotalRequests != 1 {
		t.Fatalf(fmt.Sprintf("Expected %v requests so far, but got: %v", 1, statsResponse1.TotalRequests))
	}
	if !sendSimilarWordsReq(t) {
		t.Fatalf("Failed to send a similar words request")
	}
	statsResponse2, ok := sendStatsReq(t)
	if !ok {
		t.Fatalf("Bad response for stats request")
	}
	if statsResponse2.TotalRequests != 2 {
		t.Fatalf(fmt.Sprintf("Expected %v requests so far, but got: %v", 2, statsResponse1.TotalRequests))
	}
	if statsResponse1.TotalWords != statsResponse2.TotalWords {
		t.Fatalf("Total words should not change for the lifetime of the server")
	}

}
