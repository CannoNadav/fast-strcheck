package server

import (
	"github.com/labstack/echo"
	"time"
)

func RequestDurationMiddleware(cpuTimeChannel chan<- int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip measuring time for the stats endpoint
			if c.Path() != SIMILAR_WORDS_ENDPOINT_PATH {
				return next(c)
			}

			startTime := time.Now()
			err := next(c)
			duration := time.Since(startTime)
			cpuTimeChannel <- int(duration.Nanoseconds())
			return err
		}
	}
}
