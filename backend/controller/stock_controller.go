package controller

import (
	"log"
	"net/http"
	"stock_update/models"
	"stock_update/service"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StockController handles HTTP requests related to stocks.
type StockController struct {
	service *service.StockService
}

// NewStockController creates a new StockController instance.
func NewStockController(service *service.StockService) *StockController {
	return &StockController{
		service: service,
	}
}

var (
	stocks    []models.Stock
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Stock)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mu sync.Mutex
)

// GetStocks is the handler for the endpoint to get the list of stocks.
func (controller *StockController) GetStocks(c *gin.Context) {
	stocks, err := controller.service.GetStocks()
	if err != nil {
		c.JSON(http.StatusOK, err)
	}
	c.JSON(http.StatusOK, stocks)
}

// HandleWebSocket is the handler for WebSocket connections.
func (controller *StockController) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Register client
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// Listen for new stock updates
	for {
		select {
		case stock := <-broadcast:
			// Send updated stock data to all clients
			mu.Lock()
			for client := range clients {
				if err := client.WriteJSON(stock); err != nil {
					log.Println(err)
					delete(clients, client)
					client.Close()
				}
			}
			mu.Unlock()
		}
	}
}
