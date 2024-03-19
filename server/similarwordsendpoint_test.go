package server

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"main/programargs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimilarWordsEndpointSanityTest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, SIMILAR_WORDS_ENDPOINT_PATH+"?word=bars", nil)
	rec := httptest.NewRecorder()
	cnt := e.NewContext(req, rec)

	_, wordsCatalogFactory, err := programargs.GetWordsCatalogFactory(func(s string) {
		t.Log(s)
	})

	wordsCatalog := getWordsCatalog(cnt.Logger(), wordsCatalogFactory, WORDS_CATALOG_FILE_PATH)
	err = SimilarWordsEndpoint(cnt, wordsCatalog)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	t.Log(fmt.Sprintf("Recorder, response body: %v", rec.Body))
}
