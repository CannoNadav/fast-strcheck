package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"main/wordscatalog"
	"net/http"
	"os"
	"time"
)

const (
	CHANNEL_SIZE                = 1000
	SERVER_ADDRESS              = ":8000"
	API_PREFIX                  = "/api/v1"
	SIMILAR_WORDS_ENDPOINT_PATH = API_PREFIX + "/similar"
	STATS_ENDPOINT_PATH         = API_PREFIX + "/stats"
)

type SimilarWordsService struct {
	echoServer       *echo.Echo
	statsInfoChannel chan StatsInfo
	cpuTimeChannel   chan int
}

func getWordsCatalog(logger echo.Logger, catalogFactory wordscatalog.WordsCatalogFactory, wordsCatalogFilePath string) wordscatalog.WordsCatalog {
	catalogFile, err := os.Open(wordsCatalogFilePath)
	if err != nil {
		panic("Initialization cannot proceed: Cannot open the words catalog file! - " + err.Error())
	}
	defer func() {
		if closeErr := catalogFile.Close(); closeErr != nil {
			logger.Errorf("Failed to close the words catalog file: %v", closeErr)
		}
	}()

	catalog, err := wordscatalog.ReadWordsCatalogFromFile(catalogFile, catalogFactory)
	if err != nil {
		panic(fmt.Sprintf("Initialization cannot proceed: Failed to read the application's "+
			"words catalog from file '%s', %s", wordsCatalogFilePath, err.Error()))
	}
	if catalog == nil {
		panic(fmt.Sprintf("Initialization cannot proceed with a nil words catalog"))
	}
	return catalog
}

func (similarWordsService SimilarWordsService) StopServer() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := similarWordsService.echoServer.Shutdown(ctx); err != nil {
		// not so Graceful anymore. forcefully shut down the server.
		similarWordsService.echoServer.Logger.Fatal(err)
	}
}

func (similarWordsService SimilarWordsService) StartServer() {
	StartStatsManagerGoroutine(similarWordsService.statsInfoChannel, similarWordsService.cpuTimeChannel)
	similarWordsService.echoServer.Logger.Info("Server starting")
	if err := similarWordsService.echoServer.Start(SERVER_ADDRESS); err != nil && !errors.Is(err, http.ErrServerClosed) {
		similarWordsService.echoServer.Logger.Fatal(fmt.Sprintf("Server encountered a fatal error and needs to close: %v", err))
	}
}

func InitAndStartServer(wordsCatalogFactory wordscatalog.WordsCatalogFactory, wordsCatalogFilePath string) {
	statsInfoChannel := make(chan StatsInfo, CHANNEL_SIZE)
	cpuTimeChannel := make(chan int, CHANNEL_SIZE)
	defer func() {
		close(statsInfoChannel)
		close(cpuTimeChannel)
	}()

	InitNewServer(cpuTimeChannel, statsInfoChannel, wordsCatalogFactory, wordsCatalogFilePath).StartServer()
}

func InitNewServer(cpuTimeChannel chan int, statsInfoChannel chan StatsInfo, wordsCatalogFactory wordscatalog.WordsCatalogFactory, wordsCatalogFilePath string) SimilarWordsService {
	e := echo.New()

	// Middleware
	// keep the RequestDurationMiddleware first. first middleware registered will be first
	// to receive the request amongst the chain and the last to finish, ensuring a
	// more accurate measurement of handling time.
	e.Use(RequestDurationMiddleware(cpuTimeChannel))
	//e.Use(middleware.Logger()) // this will log details on every incoming request, very verbose
	e.Use(middleware.Recover())

	wordsCatalog := getWordsCatalog(e.Logger, wordsCatalogFactory, wordsCatalogFilePath)
	numWordsInCatalog := wordsCatalog.CountWordsInCatalog()

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))

	e.GET(SIMILAR_WORDS_ENDPOINT_PATH, func(context echo.Context) error { return SimilarWordsEndpoint(context, wordsCatalog) })
	e.GET(STATS_ENDPOINT_PATH, func(context echo.Context) error {
		return GetStatsEndpoint(context, statsInfoChannel, numWordsInCatalog)
	})

	return SimilarWordsService{e, statsInfoChannel, cpuTimeChannel}
}
