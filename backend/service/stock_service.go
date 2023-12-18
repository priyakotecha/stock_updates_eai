package service

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"stock_update/models"
	"stock_update/repository"
	"time"
)

// StockService provides business logic for stocks.
type StockService struct {
	repository *repository.StockRepository
	broadcast  chan models.Stock
}

// NewStockService creates a new StockService instance.
func NewStockService(repo *repository.StockRepository, broadcast chan models.Stock) *StockService {
	return &StockService{
		repository: repo,
		broadcast:  broadcast,
	}
}

// UpdateStockPrices updates stock prices in the background.
func (s *StockService) UpdateStockPrices(stocks []models.Stock) {
	for {
		for i := range stocks {
			stock := &stocks[i]
			stock.CurrentPrice = generateRandomPrice(stock.OpenPrice)
			s.broadcast <- *stock
			time.Sleep(time.Second * time.Duration(stock.RefreshInterval))
		}
	}
}

// generateRandomPrice generates a random price change for a stock.
func generateRandomPrice(currentPrice float64) float64 {
	// For this example, used a simple random change
	change := (rand.Float64() - 0.5) * 10.0
	return currentPrice + change
}

// GetStocks get stock details.
func (s *StockService) GetStocks() ([]models.Stock, error) {
	return s.repository.GetStocks()
}

// saveStocksToFile saves the list of stocks to a JSON file.
func (s *StockService) SaveStocksToFile(stocks []models.Stock) error {
	data, err := json.MarshalIndent(stocks, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("fetched_stocks", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
