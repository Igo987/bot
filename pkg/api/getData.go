package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github/Igo87/crypt/config"
	"github/Igo87/crypt/models"
)

var ErrFailedToGet = errors.New("failed to get data")
var ErrCancelledBySignal = errors.New("canceled by signal")
var ErrFailedToCreateRequest = errors.New("failed to create request")
var ErrFailedToFetchData = errors.New("failed to fetch data")
var ErrStatusCodeNotOK = errors.New("status code not OK")
var ErrDecodeJSON = errors.New("failed to decode response")

// FetchData retrieves cryptocurrency data from the CoinMarketCap API.
func FetchData(ctx context.Context) (models.Crypto, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Cfg.GetAPIURL(), nil)
	if err != nil {
		return models.Crypto{}, ErrFailedToCreateRequest
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", config.Cfg.GetAPIKey())

	resp, err := client.Do(req)
	if err != nil {
		return models.Crypto{}, ErrFailedToFetchData
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Crypto{}, ErrStatusCodeNotOK
	}

	var data models.Crypto
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.Crypto{}, ErrDecodeJSON
	}
	return data, nil
}
