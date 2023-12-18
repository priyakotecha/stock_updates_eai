package repository

import (
	"encoding/json"
	"io/ioutil"
	"stock_update/models"
)

// StockRepository handles CRUD operations for stocks.
type StockRepository struct {
}

// NewStockRepository creates a new StockRepository instance.
func NewStockRepository() *StockRepository {
	return &StockRepository{}
}

// GetStocks returns the list of stocks by reading them from the file.
func (repo *StockRepository) GetStocks() ([]models.Stock, error) {
	data, err := ioutil.ReadFile("fetched_stocks")
	if err != nil {
		return nil, err
	}

	var stocks []models.Stock
	err = json.Unmarshal(data, &stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

// UpdateStock updates the given stock in the file.
func (repo *StockRepository) UpdateStock(updatedStock models.Stock) error {
	stocks, err := repo.GetStocks()
	if err != nil {
		return err
	}

	for i, stock := range stocks {
		if stock.ID == updatedStock.ID {
			stocks[i] = updatedStock
			break
		}
	}

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
