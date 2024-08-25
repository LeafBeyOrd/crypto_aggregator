package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Props struct {
	TokenID           string `json:"tokenId"`
	TxnHash           string `json:"txnHash"`
	ChainID           string `json:"chainId"`
	CollectionAddress string `json:"collectionAddress"`
	CurrencyAddress   string `json:"currencyAddress"`
	CurrencySymbol    string `json:"currencySymbol"`
	MarketplaceType   string `json:"marketplaceType"`
	RequestID         string `json:"requestId"`
}

type Nums struct {
	CurrencyValueDecimal string `json:"currencyValueDecimal"`
	CurrencyValueRaw     string `json:"currencyValueRaw"`
}

type Coin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type MarketData struct {
	Prices [][]float64 `json:"prices"`
}

// Cache structure for storing conversion rates
var rateCache = make(map[string]float64)
var cacheMutex = &sync.Mutex{}

// FetchCoinList fetches the list of coins and their symbols from the CoinGecko API
func FetchCoinList() map[string]string {
	url := "https://api.coingecko.com/api/v3/coins/list?include_platform=false"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", "CG-5Xf6UYaKqdcXac4BhQS8AvPc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch coin list: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var coins []Coin
	if err := json.Unmarshal(body, &coins); err != nil {
		log.Fatalf("Failed to parse coin list: %v", err)
	}

	coinMap := make(map[string]string)
	for _, coin := range coins {
		coinMap[coin.Symbol] = coin.ID
	}

	return coinMap
}

// FetchConversionRate fetches the conversion rate from a given coin to USD on a specific date
func FetchConversionRate(coinID, date string) (float64, error) {
	cacheKey := fmt.Sprintf("%s-%s", coinID, date)

	// Check if the rate is already in the cache
	cacheMutex.Lock()
	if rate, found := rateCache[cacheKey]; found {
		cacheMutex.Unlock()
		log.Printf("already exists")
		return rate, nil
	}
	cacheMutex.Unlock()

	// If not found in cache, proceed with API call
	log.Printf("Fetching conversion rate for %s on %s", coinID, date)
	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("Error parsing date: %v", err)
	}
	startTimestamp := startTime.Unix()
	endTimestamp := startTimestamp + 86400 - 1 // End of the day timestamp

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/%s/market_chart/range?vs_currency=usd&from=%d&to=%d",
		coinID, startTimestamp, endTimestamp,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", "CG-5Xf6UYaKqdcXac4BhQS8AvPc")
	log.Printf("api call  %s ", coinID)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error making API request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("Error reading response: %v", err)
	}

	// Parse the JSON response
	var data MarketData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	// Calculate the average price
	totalPrice := 0.0
	for _, price := range data.Prices {
		totalPrice += price[1]
	}
	if len(data.Prices) == 0 {
		return 0, fmt.Errorf("No price data available")
	}
	averagePrice := totalPrice / float64(len(data.Prices))

	// Store the result in the cache
	cacheMutex.Lock()
	rateCache[cacheKey] = averagePrice
	cacheMutex.Unlock()

	return averagePrice, nil
}
