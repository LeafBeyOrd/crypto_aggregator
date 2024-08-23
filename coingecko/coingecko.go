package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Coin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type MarketData struct {
	Prices [][]float64 `json:"prices"`
}

var apiKey = "CG-5Xf6UYaKqdcXac4BhQS8AvPc"

// FetchCoinList retrieves the list of coins and returns a map of symbol to ID
func FetchCoinList() map[string]string {
	url := "https://api.coingecko.com/api/v3/coins/list?include_platform=false"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", apiKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var coins []Coin
	json.Unmarshal(body, &coins)

	coinMap := make(map[string]string)
	for _, coin := range coins {
		coinMap[coin.Symbol] = coin.ID
	}

	return coinMap
}

// FetchAveragePrice fetches the average price of a coin in USD for a specific date
func FetchAveragePrice(coinID, dateStr string) float64 {
	startTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}
	startTimestamp := startTime.Unix()
	endTimestamp := startTimestamp + 86400 - 1

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/%s/market_chart/range?vs_currency=usd&from=%d&to=%d",
		coinID, startTimestamp, endTimestamp,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		return 0
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return 0
	}

	var data MarketData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return 0
	}

	// Calculate the average price
	totalPrice := 0.0
	for _, price := range data.Prices {
		totalPrice += price[1]
	}
	averagePrice := totalPrice / float64(len(data.Prices))

	return averagePrice
}
