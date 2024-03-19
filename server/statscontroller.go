package server

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sync"
)

type StatsInfo struct {
	NumRequests *int
	TimeSpent   *int
	Wg          *sync.WaitGroup
}

func StartStatsManagerGoroutine(statsInfosChannel <-chan StatsInfo, cpuTimeUpdatesChannel <-chan int) {
	go func() {
		cumulativeCPUTime := 0 // Total Time spent on requests, not including stats requests
		numRequests := 0       // Total number of requests, not including stats requests

		for {
			select {
			case cpuTimeUpdateMsg, ok := <-cpuTimeUpdatesChannel:
				if ok {
					numRequests++
					cumulativeCPUTime += cpuTimeUpdateMsg
				} else {
					cpuTimeUpdatesChannel = nil
				}
			case statsInfo, ok := <-statsInfosChannel:
				// we read the statsInfos send to this channel, but the purpose is actually to send
				// data in the opposite direction, by filling the struct sent through the channel with the necessary values.
				if ok {
					*statsInfo.NumRequests = numRequests
					*statsInfo.TimeSpent = cumulativeCPUTime
					statsInfo.Wg.Done() // signal for completion, so the goroutine on the other side can proceed
				} else {
					statsInfosChannel = nil
				}
			}

			// Break out of the loop when both channels are closed and nil
			if cpuTimeUpdatesChannel == nil && statsInfosChannel == nil {
				break
			}
		}
	}()
}

func sendStatsResponse(requestContext echo.Context, numRequests int, cumulativeCPUTimeNs int, numWordsInCatalog int) error {
	var avgProcessingTimeNs int
	if numRequests == 0 {
		avgProcessingTimeNs = 0
	} else {
		avgProcessingTimeNs = cumulativeCPUTimeNs / numRequests
	}

	statsResponse := map[string]int{
		"totalWords":          numWordsInCatalog,
		"totalRequests":       numRequests,
		"avgProcessingTimeNs": avgProcessingTimeNs,
	}
	err := requestContext.JSON(http.StatusOK, statsResponse)
	if err != nil {
		requestContext.Logger().Error(fmt.Sprintf("Failed to send stats response. Error: %s. stats payload: %v", err.Error(), statsResponse))
	}
	return err
}

func GetStatsEndpoint(context echo.Context, statsInfosChannel chan<- StatsInfo, numWordsInCatalog int) error {

	var numRequests int
	var cpuTimeSpent int
	wg := &sync.WaitGroup{}
	wg.Add(1)

	statsInfoCommand := StatsInfo{
		NumRequests: &numRequests,
		TimeSpent:   &cpuTimeSpent,
		Wg:          wg,
	}

	statsInfosChannel <- statsInfoCommand
	wg.Wait() // the goroutine spawned in sendStatsResponse will signal when the stats will be available in their memory locations

	return sendStatsResponse(context, numRequests, cpuTimeSpent, numWordsInCatalog)
}
