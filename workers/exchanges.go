package workers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ByPikod/go-crypto/helpers"
)

const API_URL = "https://api.coinbase.com/v2/exchange-rates?currency=%s"

type ExchangeRates struct {
	Currency string             `json:"currency"`
	Rates    map[string]float64 `json:"rates"`
}

type exchangeRatesParse struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

// Cached exchange rates
var (
	lastExchangeRates *ExchangeRates
)

// Returns last fetched exchange rates. Returns nil if exchange rate worker haven't been initialized.
func GetExchangeRates() *ExchangeRates {
	return lastExchangeRates
}

// Initializes exchange rate worker and keep updates the lastExchangeRates variable.
// Use GetExchangeRates() function to get exchange rates.
func InitializeExchangeRateWorker() {
	UpdateExchangeRate()
	for range time.Tick(time.Second * 30) {
		UpdateExchangeRate()
	}
}

// Fetchs and updates exchange rate
func UpdateExchangeRate() {
	exchangeRates, err := fetchExchangeRate("USD")
	if err != nil {
		helpers.LogError("Failed to fetch exchange rate: " + err.Error())
	}

	lastExchangeRates = exchangeRates
}

// Fetchs API and retrieves exchange rates data in the form of ExchangeRates struct.
func fetchExchangeRate(currency string) (*ExchangeRates, error) {
	res, err := http.Get(fmt.Sprintf(API_URL, currency))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	result := struct {
		Data exchangeRatesParse
	}{}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return parseExchanges(&result.Data)
}

// Parse string and convert it to float64
func parseExchanges(exchangeRates *exchangeRatesParse) (*ExchangeRates, error) {
	// parse floats
	parsed := ExchangeRates{Rates: map[string]float64{}}
	parsed.Currency = exchangeRates.Currency
	for currency, rateToParse := range exchangeRates.Rates {
		rate, err := strconv.ParseFloat(rateToParse, 64)
		if err != nil {
			return nil, err
		}
		parsed.Rates[currency] = rate
	}
	// return
	return &parsed, nil
}
