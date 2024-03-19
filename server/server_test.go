// server/server_test.go

package server

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/labstack/echo"
	"main/programargs"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

const WORDS_CATALOG_FILE_PATH = "../words_clean.txt"

var testResult string

var testStartTime time.Time
var testEndTime time.Time

type CSVData struct {
	t      time.Time
	memUse float64
}

func startMemoryProfiling(t *testing.T, csvFileWriteWaitGroup *sync.WaitGroup, stopMemoryProfilingChannel <-chan bool, sleepInterval time.Duration, wordscatalogAlgName string, testResultsFolderPath string) {
	var m runtime.MemStats
	bufferSize := 100
	csvData := make([]CSVData, bufferSize)
	sumMemory := float64(0)
	numDataPoints := 0
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var CSV_FILE_PATH = testResultsFolderPath + wordscatalogAlgName + "_memory_usage.csv"
		defer func() { csvFileWriteWaitGroup.Done() }() // will be called last amongst defers, after Write.Flush() finished

		// Save to CSV
		csvFile, err := os.Create(CSV_FILE_PATH)
		if err != nil {
			t.Errorf("Could not create CSV file '" + CSV_FILE_PATH + "'. " + err.Error())
			return
		}
		defer func() {
			if err = csvFile.Close(); err != nil {
				t.Errorf("Could not close the CSV file '" + CSV_FILE_PATH + "'. " + err.Error())
			}
		}()

		writer := csv.NewWriter(csvFile)
		// Write header and data to the CSV file
		if err = writer.Write([]string{"Time", "Memory (MB)"}); err != nil {
			t.Errorf("Could not write memory CSV file header: " + err.Error())
			return
		}
		writer.Flush()

		debug.FreeOSMemory() // giving a clean start to memory stats
		wg.Done()

		for {
			select {
			case <-stopMemoryProfilingChannel:
				flushCSVSliceToFile(t, writer, &csvData, &sumMemory, &numDataPoints)
				var avgMemoryConsumption float64
				if numDataPoints == 0 {
					avgMemoryConsumption = 0
				} else {
					avgMemoryConsumption = sumMemory / float64(numDataPoints)
				}

				switch testResult {
				case "Pass":
					testResult = "Pass"
				case "Fail":
					testResult = "Fail"
				case "":
					testResult = "Unknown(test status not received from client)"
				default:
					testResult = fmt.Sprintf("Unrecognized option received from client('%v')", testResult)
				}

				t.Logf("\n\nTest finished. Results(Summary):\n---------------> Status: %v\n---------------> Words catalog type: %v\n---------------> Test time: %v seconds\n---------------> Average memory consumption during test: %v(MB)\n\n", testResult, wordscatalogAlgName, testEndTime.Sub(testStartTime).Seconds(), avgMemoryConsumption)
				return
			default:
				runtime.ReadMemStats(&m)
				alloc := float64(m.Alloc) / (1024 * 1024) // Convert bytes to MB
				// Append data for CSV
				csvData = append(csvData, CSVData{t: time.Now(), memUse: alloc})
				if len(csvData) >= bufferSize {
					flushCSVSliceToFile(t, writer, &csvData, &sumMemory, &numDataPoints)
				}
				time.Sleep(sleepInterval)
			}
		}
	}()

	wg.Wait()
}

func flushCSVSliceToFile(t *testing.T, writer *csv.Writer, csvData *[]CSVData, sumMemory *float64, numDatapoints *int) {
	var csvDataStrings [][]string
	for _, dataPoint := range *csvData {
		if testStartTime.IsZero() || dataPoint.t.Before(testStartTime) || (!testEndTime.IsZero() && dataPoint.t.After(testEndTime)) {
			continue
		}
		csvDataStrings = append(csvDataStrings, []string{fmt.Sprintf("%v", dataPoint.t.Sub(testStartTime).Nanoseconds()), fmt.Sprintf("%.3f", dataPoint.memUse)})
		*sumMemory += dataPoint.memUse
		*numDatapoints++

	}
	if err := writer.WriteAll(csvDataStrings); err != nil {
		t.Errorf("Could not write memory CSV file data: " + err.Error())
	}
	*csvData = (*csvData)[:0] // Reset the length of the slice to 0
}

func TestServerWithProfiling(t *testing.T) {
	// This test is expected to be run in conjunction with one of the test clients found in the test_client package

	flag.Parse()

	testResultsFolderPath, err := filepath.Abs("../server_test_results/")
	testResultsFolderPath = testResultsFolderPath + "/"

	if err != nil {
		t.Fatal("Could not obtain path for test results folder:", err)
	}
	err = os.MkdirAll(testResultsFolderPath, 0755)
	if err != nil {
		t.Fatal("Error creating directory for test results:", err)
	}

	t.Logf("This test must be run in conjunction with one of the test clients. for possible clients see tests in the 'client' directory.")
	wordscatalogAlgName, wordsCatalogFactory, err := programargs.GetWordsCatalogFactory(func(s string) { t.Log(s) })
	if err != nil {
		t.Fatal("Problem obtaining a words catalog factory: ", err.Error())
	}
	cpuPProfFilePath := testResultsFolderPath + wordscatalogAlgName + "_cpu.pprof"
	cpuFile, err := os.Create(cpuPProfFilePath)
	if err != nil {
		t.Fatalf("could not create CPU profile: %v", err)
	}

	defer func() {
		if err = cpuFile.Close(); err != nil {
			t.Fatalf("Failed to close resource: %v", err)
		}
	}()

	var shutdownEndpointWaitGroup sync.WaitGroup
	var csvFileWriteWaitGroup sync.WaitGroup
	stopMemoryProfilingChannel := make(chan bool)
	statsInfoChannel := make(chan StatsInfo, CHANNEL_SIZE)
	cpuTimeChannel := make(chan int, CHANNEL_SIZE)

	similarWordsService := InitNewServer(cpuTimeChannel, statsInfoChannel, wordsCatalogFactory, WORDS_CATALOG_FILE_PATH)
	//similarWordsService.echoServer.Logger.SetOutput(io.Discard) // discarding logger output for performance reasons during benchmarking

	// add an endpoint for the test client to signal when he finished sending test inputs
	var testFinishGuard sync.Once
	shutdownEndpointWaitGroup.Add(1) // keeps the server alive until we get the signal to shutdown from the client
	similarWordsService.echoServer.GET("/stop-clock-and-shutdown", func(c echo.Context) error {
		testFinishGuard.Do(func() {
			t.Log("Client signaled to finish test, stopping test timer.")
			testResult = c.QueryParam("testResult")
			testEndTime = time.Now()
			stopMemoryProfilingChannel <- true
			shutdownEndpointWaitGroup.Done()
		})

		return nil
	})

	var startTimerOnce sync.Once
	similarWordsService.echoServer.GET("/start-clock", func(c echo.Context) error {
		startTimerOnce.Do(func() {
			t.Log("Client signaled to start test, starting test timer.")
			runtime.GC() // start the test from with clean gc
			testStartTime = time.Now()
		})
		return c.String(200, "ok")
	})

	go similarWordsService.StartServer()
	defer similarWordsService.StopServer()

	// Start CPU profiling
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		t.Fatalf(fmt.Sprintf("could not start cpu profiling: %v", err))
	}
	defer pprof.StopCPUProfile()
	defer func() {
		close(stopMemoryProfilingChannel)
		close(statsInfoChannel)
		close(cpuTimeChannel)
		csvFileWriteWaitGroup.Wait() // don't close the main goroutine before the csv file was written
	}()

	csvFileWriteWaitGroup.Add(1)

	startMemoryProfiling(t, &csvFileWriteWaitGroup, stopMemoryProfilingChannel, 1*time.Millisecond, wordscatalogAlgName, testResultsFolderPath)
	// only open the prob-server-ready endpoint after all profiling and server have started
	similarWordsService.echoServer.GET("/prob-server-ready", func(c echo.Context) error {
		return c.String(200, "ok")
	})
	shutdownEndpointWaitGroup.Wait()

}
