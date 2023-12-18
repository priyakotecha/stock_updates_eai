package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"stock_update/controller"
	"stock_update/models"
	"stock_update/repository"
	"stock_update/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch stocks from Polygon API and store them in a file
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		log.Fatal("Polygon API key is missing. Set POLYGON_API_KEY environment variable.")
	}
	stocks, err := fetchTopStocks(apiKey)
	log.Println("Stocks Fetched ", stocks)
	// Create repository, service, and controller
	stockRepository := repository.NewStockRepository()
	broadcast := make(chan models.Stock)
	stockService := service.NewStockService(stockRepository, broadcast)
	err = stockService.SaveStocksToFile(stocks)
	if err != nil {
		log.Fatal("Error adding stocks to File")
	}
	stockController := controller.NewStockController(stockService)
	// Start updating stock prices in the background
	// stockRepository.UpdateStock(stocks)
	go stockService.UpdateStockPrices(stocks)

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	api := router.Group("/api")
	{
		api.GET("/stocks", stockController.GetStocks)
		api.GET("/ws", stockController.HandleWebSocket)
	}

	// Start the server
	port := os.Getenv("PORT")
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(router.Run(":" + port))
}

// fetchTopStocks fetches the list of top 20 stocks along with their open prices from the Polygon API.
func fetchTopStocks(apiKey string) ([]models.Stock, error) {
	url := "https://api.polygon.io/v3/reference/tickers?active=true&order=asc&limit=20&sort=ticker&apiKey" + apiKey

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch top stocks: %s", response.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	var stocks []models.Stock
	if tickers, ok := result["tickers"].([]interface{}); ok {
		for _, ticker := range tickers {
			if tickerMap, ok := ticker.(map[string]interface{}); ok {
				symbol := tickerMap["symbol"].(string)
				name := tickerMap["name"].(string)
				market := tickerMap["market"].(string)

				// You may need to make additional requests to get the open price or other details
				openPrice, err := getOpenPrice(apiKey, symbol)
				if err != nil {
					log.Printf("Error fetching open price for %s: %v", symbol, err)
					continue
				}

				stock := models.Stock{
					Symbol:    symbol,
					Name:      name,
					Market:    market,
					OpenPrice: openPrice,
				}
				stocks = append(stocks, stock)

			}
		}
	}

	return stocks, nil
}

// getOpenPrice fetches the open price for a given stock symbol from the Polygon API.
func getOpenPrice(apiKey, symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?apiKey=%s", symbol, apiKey)

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	if results, ok := result["results"].([]interface{}); ok {
		if len(results) > 0 {
			if firstResult, ok := results[0].(map[string]interface{}); ok {
				openPrice, _ := firstResult["o"].(float64)
				return openPrice, nil
			}
		}
	}

	return 0, fmt.Errorf("Open price not found for %s", symbol)
}
