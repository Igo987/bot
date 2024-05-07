package models

import (
	"time"
)

type Crypto struct {
	Data struct {
		Bitcoin struct {
			Name        string    `json:"name"`
			Symbol      string    `json:"symbol"`
			LastUpdated time.Time `json:"last_updated"`
			Quote       struct {
				Rub struct {
					Price            float64 `json:"price"`
					PercentChange1H  float64 `json:"percent_change_1h"`
					PercentChange24H float64 `json:"percent_change_24h"`
					PercentChange7D  float64 `json:"percent_change_7d"`
					PercentChange30D float64 `json:"percent_change_30d"`
				} `json:"RUB"`
			} `json:"quote"`
		} `json:"1"`
		Ethereum struct {
			Name        string    `json:"name"`
			Symbol      string    `json:"symbol"`
			LastUpdated time.Time `json:"last_updated"`
			Quote       struct {
				Rub struct {
					Price            float64 `json:"price"`
					PercentChange1H  float64 `json:"percent_change_1h"`
					PercentChange24H float64 `json:"percent_change_24h"`
					PercentChange7D  float64 `json:"percent_change_7d"`
					PercentChange30D float64 `json:"percent_change_30d"`
				} `json:"RUB"`
			} `json:"quote"`
		} `json:"1027"`
	} `json:"data"`
}

type Extremes struct {
	Name    string  `json:"name"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Percent float64 `json:"percent"`
}

type Currencies []Extremes

type CurrencyValue struct {
	Price            float64 `json:"price"`
	PercentChange1H  float64 `json:"percent_change_1h"`
	PercentChange24H float64 `json:"percent_change_24h"`
	PercentChange7D  float64 `json:"percent_change_7d"`
	PercentChange30D float64 `json:"percent_change_30d"`
}
