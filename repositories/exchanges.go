package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ByPikod/go-crypto/helpers"
)

type (
	ExchangeRates struct {
		Currency string             `json:"currency"`
		Rates    map[string]float64 `json:"rates"`
	}
	ExchangesRepository struct {
		lastExchangeData *ExchangeRates
	}
)

const API_URL = "https://api.coinbase.com/v2/exchange-rates?currency=%s"

// New exchanges repository
func NewExchangesRepository() *ExchangesRepository {
	res := &ExchangesRepository{
		lastExchangeData: nil,
	}
	// Start fetching data
	go res.initializeExchangeRateUpdater()
	return res
}

// Fetchs API and retrieves exchange rates data in the form of ExchangeRates struct.
func (repo *ExchangesRepository) fetchExchangeRate(currency string) (*ExchangeRates, error) {

	// Fetch
	res, err := http.Get(fmt.Sprintf(API_URL, currency))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Unpack/Decode
	result := struct {
		Data struct {
			Currency string            `json:"currency"`
			Rates    map[string]string `json:"rates"`
		}
	}{}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Parse
	parsedRates, err := repo.parseFloats(result.Data.Rates)
	if err != nil {
		return nil, err
	}

	return &ExchangeRates{
		Currency: result.Data.Currency,
		Rates:    parsedRates,
	}, nil
}

// Parse string and convert it to float64
func (repo *ExchangesRepository) parseFloats(exchangeRates map[string]string) (map[string]float64, error) {
	// result array
	parsed := map[string]float64{}
	// loop it and parse one by one
	for currency, rateToParse := range exchangeRates {
		rate, err := strconv.ParseFloat(rateToParse, 64)
		if err != nil {
			return nil, err
		}
		parsed[currency] = rate
	}
	// return
	return parsed, nil
}

// Fetchs and updates exchange rate
func (repo *ExchangesRepository) updateExchangeRate() {
	exchangeRates, err := repo.fetchExchangeRate("USD")
	if err != nil {
		helpers.LogError("Failed to fetch exchange rate: " + err.Error())
	}

	repo.lastExchangeData = exchangeRates
}

// Initializes exchange rate worker and keep updates the lastExchangeRates variable.
// Use GetExchangeRates() function to get exchange rates.
func (repo *ExchangesRepository) initializeExchangeRateUpdater() {
	repo.updateExchangeRate()
	for range time.Tick(time.Second * 30) {
		repo.updateExchangeRate()
	}
}

// Returns last fetched exchange rates. Returns nil if exchange rate worker haven't been initialized.
func (repo *ExchangesRepository) GetExchangeRates() *ExchangeRates {
	return repo.lastExchangeData
}
