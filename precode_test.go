package main

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandlerWhenOk(t *testing.T) {
	req := newTestRequest("GET", "/cafe?count=4&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()

	require.NotNil(t, req, "Error creating request")
	require.Equal(t, http.StatusOK, response.StatusCode, "Unexpected status code")
	require.NotEmpty(t, response.Body, "Response body is empty")
}

func TestMainHandlerWhereIsTheWrongCity(t *testing.T) {
	req := newTestRequest("GET", "/cafe?count=4&city=Ryazan", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Unexpected status code")

	body, _ := io.ReadAll(response.Body)
	require.Contains(t, string(body), "wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := newTestRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()

	body, _ := io.ReadAll(response.Body)
	list := strings.Split(string(body), ",")

	assert.Len(t, list, totalCount, "Unexpected number of cafes")
}

func newTestRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	return req
}
