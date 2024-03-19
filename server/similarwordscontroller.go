package server

import (
	"github.com/labstack/echo"
	"main/wordscatalog"
	"net/http"
	"strings"
)

var errWordParamIsMissingResponse map[string]string

func init() {
	errWordParamIsMissingResponse = map[string]string{
		"error": "word parameter is missing",
	}
}

func SimilarWordsEndpoint(c echo.Context, wordsCatalog wordscatalog.WordsCatalog) error {
	// this endpoint doesn't write to the cpu time channel itself(as it doesn't know how much time the
	// request took). writing this value is done for it in TimeMeasuringMiddleware
	word := c.QueryParam("word")
	if word == "" {
		c.Logger().Error("SimilarWordsRequest with a missing word parameter")
		return c.JSON(http.StatusBadRequest, errWordParamIsMissingResponse)
	}

	equivalentWords := wordsCatalog.FindEquivalentWords(strings.ToLower(word))
	if equivalentWords == nil || len(equivalentWords) == 0 {
		return c.JSON(http.StatusOK, map[string][]string{
			"similar": {},
		})
	} else {
		return c.JSON(http.StatusOK, map[string][]string{
			"similar": equivalentWords,
		})
	}

}
