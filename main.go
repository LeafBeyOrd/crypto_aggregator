package main

import (
	"crypto-aggregator/coingecko"
	"crypto-aggregator/utils"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Props struct {
	CurrencySymbol string `json:"currencySymbol"`
}

type Nums struct {
	CurrencyValueDecimal string `json:"currencyValueDecimal"`
}

func main() {
	// Fetch the coin list and create the coinMap
	coinMap := coingecko.FetchCoinList()

	// Open the CSV file
	f, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	// Create a CSV reader
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // Allows variable number of fields per record

	// Skip the header row
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}

	// Aggregate data
	data := make(map[string]map[string]float64)

	// Process each row
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read record: %v", err)
		}

		ts := record[1]
		projectID := record[3]
		propsStr := record[14]
		numsStr := record[15]

		var props Props
		if err := json.Unmarshal([]byte(propsStr), &props); err != nil {
			log.Printf("Failed to parse props: %v; raw input: %s", err, propsStr)
			continue
		}

		var nums Nums
		if err := json.Unmarshal([]byte(numsStr), &nums); err != nil {
			log.Printf("Failed to parse nums: %v; raw input: %s", err, numsStr)
			continue
		}

		// Parse date
		date, err := time.Parse("2006-01-02 15:04:05.000", ts)
		if err != nil {
			log.Printf("Failed to parse timestamp: %v; raw input: %s", err, ts)
			continue
		}
		dateStr := date.Format("2006-01-02")

		// Fetch the coin ID
		coinID := coinMap[strings.ToLower(props.CurrencySymbol)]
		if coinID == "" {
			log.Printf("Coin ID not found for symbol: %s", props.CurrencySymbol)
			continue
		}

		// Fetch the average price in USD for that date
		averagePrice := coingecko.FetchAveragePrice(coinID, dateStr)
		if averagePrice == 0 {
			log.Printf("Failed to fetch price for %s on %s", coinID, dateStr)
			continue
		}

		// Convert the value to float
		currencyValue, err := strconv.ParseFloat(nums.CurrencyValueDecimal, 64)
		if err != nil {
			log.Printf("Failed to parse currency value: %v; raw input: %s", err, nums.CurrencyValueDecimal)
			continue
		}

		// Calculate total volume in USD
		totalVolume := currencyValue * averagePrice

		// Aggregate the data
		if _, exists := data[dateStr]; !exists {
			data[dateStr] = make(map[string]float64)
		}
		data[dateStr][projectID] += totalVolume
	}

	// Write to output.csv
	utils.WriteOutput(data)
}
