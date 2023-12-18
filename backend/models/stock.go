package models

type Stock struct {
	ID              string  `json:"id"`
	Symbol          string  `json:"symbol"`
	OpenPrice       float64 `json:"openPrice"`
	Name            string  `json:"name"`
	Market          string  `json:"market"`
	CurrentPrice    float64 `json:"currentPrice"`
	RefreshInterval int     `json:"refreshInterval"`
}
