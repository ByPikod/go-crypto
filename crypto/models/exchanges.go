package models

type ExchangeRates struct {
	Currency string             `json:"currency"`
	Rates    map[string]float64 `json:"rates"`
}
