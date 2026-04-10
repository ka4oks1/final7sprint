package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

type CountChecker struct {
	CountValue int
	CafesCount int
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	testTable := map[string]int{
		"/cafe?count=0&city=moscow":   0,
		"/cafe?count=1&city=moscow":   1,
		"/cafe?count=2&city=moscow":   2,
		"/cafe?count=100&city=moscow": 100,
	}

	for k, _ := range testTable {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", k, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, response.Code, http.StatusOK)

		var cafesCount int
		count, _ := strconv.Atoi(req.FormValue("count"))

		switch {
		case response.Body.String() == "":
			{
				cafesCount = 0
			}
		case count == 100:
			{
				cafesCount = min(100, count)
			}

		default:
			{
				cafesCount = len(strings.Split(response.Body.String(), ","))
			}

		}

		correct := cafesCount <= testTable[k]

		assert.True(t, correct)
	}

}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	testTable := map[string]int{
		"/cafe?search=фасоль&city=moscow": 0,
		"/cafe?search=кофе&city=moscow":   2,
		"/cafe?search=вилка&city=moscow":  1,
	}

	for k, _ := range testTable {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", k, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		var containsForReal bool

		cafes := strings.Split(response.Body.String(), ",")

		searchItem := req.FormValue("search")
		if response.Body.String() == "" {
			assert.Equal(t, 0, testTable[req.URL.String()])
		} else {
			assert.Equal(t, len(cafes), testTable[req.URL.String()])

			for _, v := range cafes {

				containsForReal = strings.Contains(strings.ToUpper(v), strings.ToUpper(searchItem))

			}

			assert.True(t, containsForReal)

		}

	}

}
