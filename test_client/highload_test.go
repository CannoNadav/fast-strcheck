package test_client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func stressTest(t *testing.T, allWordsInCatalog []string) bool {
	const similarWordsGoroutines = 1000
	const queriesPerGoroutine = 1000
	const statsEndpointGoroutines = 50
	var numSimilarWordsRequestsFinished atomic.Int32
	var numStatsRequestsFinished atomic.Int32
	var wg sync.WaitGroup
	var didTestFail atomic.Bool // completely unnecessary to be atomic atm since there are only concurrent writes are of the same value(true) but better still practice..
	didTestFail.Store(false)

	// check stats response when we still know what the state should be, before similar words were sent
	statsResponse1, isFail1 := SendStatsReq(t, &numStatsRequestsFinished, len(allWordsInCatalog))
	statsResponse2, isFail2 := SendStatsReq(t, &numStatsRequestsFinished, len(allWordsInCatalog))
	if int(numStatsRequestsFinished.Load()) == 2 {
		numStatsRequestsFinished.Store(0) // reset the counter for clean results in rest of the test
	} else {
		didTestFail.Store(true)
		t.Error(fmt.Sprintf("Something is off, expected 2 stats requests so far, got: %v", int(numStatsRequestsFinished.Load())))
	}

	if isFail1 || isFail2 || statsResponse1.TotalRequests != 0 || statsResponse2.TotalRequests != 0 || statsResponse1.AvgProcessingTimeNs != 0 || statsResponse2.AvgProcessingTimeNs != 0 {
		t.Error(fmt.Sprintf("Bad values before test began: %v, %v", statsResponse1.String(), statsResponse2.String()))
	}

	for i := 0; i < similarWordsGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < queriesPerGoroutine; j++ {
				success := sendSimilarWordsRequestWithoutUnmarshalling(t, allWordsInCatalog, &numSimilarWordsRequestsFinished)
				if !success {
					didTestFail.Store(true)
				}
			}
		}()
	}

	for i := 0; i < statsEndpointGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < queriesPerGoroutine; j++ {
				_, failTest := SendStatsReq(t, &numStatsRequestsFinished, len(allWordsInCatalog))
				if failTest {
					didTestFail.Store(true)
				}
			}
		}()

	}

	wg.Wait()

	totalNumSimilarWordsRequestsMade := int(numSimilarWordsRequestsFinished.Load())
	expected := similarWordsGoroutines * queriesPerGoroutine
	if totalNumSimilarWordsRequestsMade != expected {
		t.Error(fmt.Sprintf("Unexpected amount of similar words requests requests finished. expected: %v, actual: %v ", expected, totalNumSimilarWordsRequestsMade))
		didTestFail.Store(true)
	}
	totalNumStatsRequestsMade := int(numStatsRequestsFinished.Load())
	expected = statsEndpointGoroutines * queriesPerGoroutine
	if totalNumStatsRequestsMade != expected {
		t.Error(fmt.Sprintf("Unexpected amount of stats requests were made. expected: %v, actual: %v ", expected, totalNumStatsRequestsMade))
		didTestFail.Store(true)
	}

	t.Log(fmt.Sprintf("Client: Finished sending(and validating) %v similar words requests and %v stats requests", totalNumSimilarWordsRequestsMade, totalNumStatsRequestsMade))

	// doing two more stats request, since now we know exactly what response to expect(since there is no concurrency involved)
	// so we can validate the totalNumSimilarWordsRequestsMade field of the result
	statsResponse1, isFail1 = SendStatsReq(t, &numStatsRequestsFinished, len(allWordsInCatalog))
	statsResponse2, isFail2 = SendStatsReq(t, &numStatsRequestsFinished, len(allWordsInCatalog))

	// most fields are verified inside sendStatsReq
	if isFail1 || isFail2 {
		didTestFail.Store(true)
	}
	if statsResponse1.AvgProcessingTimeNs != statsResponse2.AvgProcessingTimeNs {
		t.Error(fmt.Sprintf("Expected same avg processing time(counts only similar words requests) but %v != %v", statsResponse1.AvgProcessingTimeNs, statsResponse2.AvgProcessingTimeNs))
		didTestFail.Store(true)
	}
	if statsResponse1.TotalRequests != totalNumSimilarWordsRequestsMade || statsResponse2.TotalRequests != totalNumSimilarWordsRequestsMade {
		t.Error(fmt.Sprintf("Expected number of requests in both to be %v but was %v, %v", totalNumSimilarWordsRequestsMade, statsResponse1.TotalRequests, statsResponse2.TotalRequests))
		didTestFail.Store(true)
	}

	return didTestFail.Load()
}

func TestHighLoadTestClient(t *testing.T) {
	RunTest(t, "high load test", stressTest)
}
