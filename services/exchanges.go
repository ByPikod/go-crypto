package services

import (
	"time"

	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/gofiber/contrib/websocket"
)

type (
	ExchangesService struct {
		repository repositories.IExchangesRepository
		clients    map[*websocket.Conn]chan *models.ExchangeRates
	}
)

// Broadcast the last exchange data to all the clients with an interval.
func (service *ExchangesService) wsExchangeBroadcaster() {
	// Wait
	for range time.Tick(5 * time.Second) {
		lastExchangeData := service.repository.GetExchangeRates()
		// Broadcast the last exchange data to all the clients connected.
		for _, ch := range service.clients {
			ch <- lastExchangeData
		}
	}
}

// This will create a new exchange service and start to broadcast clients that added to service.
func NewExchangeService(repository *repositories.ExchangesRepository) *ExchangesService {
	service := &ExchangesService{
		repository: repository,
		clients:    make(map[*websocket.Conn]chan *models.ExchangeRates),
	}
	// start broadcasting
	go service.wsExchangeBroadcaster()
	return service
}

// Returns last fetched exchange rates. Returns nil if exchange rate worker haven't been initialized.
func (service *ExchangesService) GetExchangeRates() (exchangeRates *models.ExchangeRates) {
	exchangeRates = service.repository.GetExchangeRates()
	return
}

// Returns the rate with "comma ok"
func (service *ExchangesService) GetCurrency(currency string) (rate float64, ok bool) {
	rate, ok = service.GetExchangeRates().Rates[currency]
	return
}

// Add websocket client to the listeners
func (service *ExchangesService) AddClient(client *websocket.Conn) chan *models.ExchangeRates {
	// Create a channel to receive broadcasts
	ch := make(chan *models.ExchangeRates)
	// Add client to the listeners
	service.clients[client] = ch
	return ch
}

// Remove websocket client from the listeners
func (service *ExchangesService) RemoveClient(client *websocket.Conn) {
	close(service.clients[client])
	delete(service.clients, client)
}
