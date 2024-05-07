package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	testServer := httptest.NewServer(handler)

	defer testServer.Close()
	req := httptest.NewRequest(http.MethodGet, testServer.URL, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler(t *testing.T) {
	t.Parallel()
	// Prepare the test data once, outside of the test server
	testData := `[{"name":"Bitcoin","min":6153047.733307837,"max":6190778.29529178,"percent":0.6132011910081395},
	{"name":"Ethereum","min":297529.2481910594,"max":301686.11320297676,"percent":1.3971281940147358}]`
	// Prepare the handler once, outside of the test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		switch r.URL.Path {
		case "/api/coins":
			data = []byte(testData)
		case "/api/coins/Bitcoin":
			data = []byte(`{"name":"Bitcoin","min":6153047.733307837,"max":6190778.29529178,"percent":0.6132011910081395}`)
		default:
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)

		_, err := w.Write(data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
	testServer := httptest.NewServer(handler)

	defer testServer.Close()
	// Test cases
	testCases := []struct {
		path   string
		status int
		data   []byte
	}{
		{"/api/coins", http.StatusOK, []byte(testData)},
		{"/api/coins/Bitcoin", http.StatusOK, []byte(`{"name":"Bitcoin","min":6153047.733307837,
		"max":6190778.29529178,"percent":0.6132011910081395}`)},
	}
	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, testServer.URL+testCase.path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != testCase.status {
			t.Errorf("for %s, want %d, got %d", testCase.path, testCase.status, rr.Code)
		}
		if rr.Body.String() != string(testCase.data) {
			t.Errorf("for %s, want %s, got %s", testCase.path, testCase.data, rr.Body.String())
		}
	}
}
