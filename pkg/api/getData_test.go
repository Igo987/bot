package api_test

import (
	"context"
	"github/Igo87/crypt/config"
	"github/Igo87/crypt/models"
	"github/Igo87/crypt/pkg/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchData_Success(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{}`))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))

	defer server.Close()

	config.Cfg.API.URL = server.URL
	data, err := api.FetchData(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, data)
}

func TestFetchData_NewRequestError(t *testing.T) {
	t.Parallel()
	// Simulate an error when creating a new request
	// config.Cfg.SetApiUrl("invalid-url") // Set an invalid URL to trigger an error
	config.Cfg.API.URL = "invalid-url"
	ctx := context.Background()
	data, err := api.FetchData(ctx)
	require.Error(t, err)
	assert.Equal(t, models.Crypto{}, data)
}

func TestFetchData_RequestError(t *testing.T) {
	t.Parallel()
	// Simulate an error when making the request
	server := httptest.NewServer(nil) // Create a server that always returns an error
	defer server.Close()
	config.Cfg.API.URL = server.URL
	ctx := context.Background()
	data, err := api.FetchData(ctx)
	require.Error(t, err)
	assert.Equal(t, models.Crypto{}.Data, data)
}

func TestFetchData_UnexpectedStatusCode(t *testing.T) {
	t.Parallel()
	// Create a test server that returns a non-OK status code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound) // Simulating a 404 response
	}))
	defer server.Close()

	config.Cfg.API.URL = server.URL

	ctx := context.Background()
	data, err := api.FetchData(ctx)
	require.Error(t, err)
	assert.Equal(t, models.Crypto{}.Data, data)
}

func TestFetchData_DecodeError(t *testing.T) {
	t.Parallel()
	// Create a test server that returns invalid JSON data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`invalid-json-response`))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	config.Cfg.API.URL = server.URL

	ctx := context.Background()
	data, err := api.FetchData(ctx)
	require.Error(t, err)
	assert.Equal(t, models.Crypto{}, data)
}
