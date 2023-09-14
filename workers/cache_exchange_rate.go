package workers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ByPikod/go-crypto/helpers"
)

const API_URL = "https://api.coinbase.com/v2/exchange-rates?currency=%s"

type exchangeRatesData struct {
	Data ExchangeRates `json:"data"`
}

type ExchangeRates struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

// Cached exchange rates
var lastExchangeRates *ExchangeRates

// Returns last fetched exchange rates. Returns nil if exchange rate worker haven't been initialized.
func GetExchangeRates() *ExchangeRates {
	return lastExchangeRates
}

// Initializes exchange rate worker and keep updates the lastExchangeRates variable.
// Use GetExchangeRates() function to get exchange rates.
func InitializeExchangeRateWorker() {

	for range time.Tick(time.Second * 5) {
		exchangeRates, err := fetchExchangeRate("USD")
		if err != nil {
			helpers.LogError("Failed to fetch exchange rate: " + err.Error())
		}

		lastExchangeRates = exchangeRates
	}

}

// Fetchs API and retrieves exchange rates data in the form of ExchangeRates struct.
func fetchExchangeRate(currency string) (*ExchangeRates, error) {
	res, err := http.Get(fmt.Sprintf(API_URL, currency))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var result exchangeRatesData

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}
