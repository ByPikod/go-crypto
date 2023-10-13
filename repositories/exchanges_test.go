package repositories

import "github.com/ByPikod/go-crypto/tree/crypto/models"

type (
	TestExchangesRepository struct {
	}
)

// New exchange repository but for testing.
func NewTestExchangesRepository() *TestExchangesRepository {
	return &TestExchangesRepository{}
}

// Returns a clone of the exchange rates for testing.
// Exchange rate data for testing only contains BTC and ETH prices.
func (repo *TestExchangesRepository) GetExchangeRates() *models.ExchangeRates {
	return &models.ExchangeRates{
		Currency: "USD",
		Rates: map[string]float64{
			"BTC": 0.00003595314945,
			"ETH": 0.0006168707994954,
		},
	}
}
