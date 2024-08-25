package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"crypto-aggregator/coingecko"
	"crypto-aggregator/utils"
)

type AggregatedData struct {
	ProjectID           string
	Date                string
	NumberOfTransactions int
	TotalVolumeUSD      float64
}

func main() {
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

	// Initialize a map to store aggregated data
	aggregatedData := make(map[string]*AggregatedData)

	// Fetch the coin list for mapping symbols to IDs
	coinList := coingecko.FetchCoinList()

	// Process each row
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		ts := record[1]           // Assuming this is the column for timestamp
		projectID := record[3]    // Assuming this is the column for project ID
		propsStr := record[14]    // Assuming this is the column for props JSON
		numsStr := record[15]     // Assuming this is the column for nums JSON

		// Parse the props JSON
		var props coingecko.Props
		if err := json.Unmarshal([]byte(propsStr), &props); err != nil {
			log.Printf("Failed to parse props: %v; raw input: %s", err, propsStr)
			continue
		}

		// Parse the nums JSON
		var nums coingecko.Nums
		if err := json.Unmarshal([]byte(numsStr), &nums); err != nil {
			log.Printf("Failed to parse nums: %v; raw input: %s", err, numsStr)
			continue
		}

		// Get the date part from the timestamp
		date := utils.ParseDate(ts)

		// Generate a unique key for the date and project ID
		key := fmt.Sprintf("%s_%s", date, projectID)

		// Convert the currency to USD
		coinID := coinList[strings.ToLower(props.CurrencySymbol)]
		if coinID == "" {
			log.Printf("Currency symbol %s not found in coin list", props.CurrencySymbol)
			continue
		}
		log.Printf("FetchConversionRate:")
		conversionRate, err := coingecko.FetchConversionRate(coinID, date)
		if err != nil {
			log.Printf("Failed to fetch conversion rate for %s: %v", coinID, err)
			continue
		}

		// Calculate the volume in USD
		volumeUSD, err := strconv.ParseFloat(nums.CurrencyValueDecimal, 64)
		if err != nil {
			log.Printf("Failed to parse currency value: %v", err)
			continue
		}
		volumeUSD *= conversionRate

		// Aggregate the data
		if data, exists := aggregatedData[key]; exists {
			data.NumberOfTransactions++
			data.TotalVolumeUSD += volumeUSD
		} else {
			aggregatedData[key] = &AggregatedData{
				ProjectID:           projectID,
				Date:                date,
				NumberOfTransactions: 1,
				TotalVolumeUSD:      volumeUSD,
			}
		}
	}

	// Open the output CSV file
	outputFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write the header
	writer.Write([]string{"date", "project_id", "number_of_transactions", "total_volume_usd"})

	// Write the aggregated data
	for _, data := range aggregatedData {
		writer.Write([]string{
			data.Date,
			data.ProjectID,
			fmt.Sprintf("%d", data.NumberOfTransactions),
			fmt.Sprintf("%.2f", data.TotalVolumeUSD),
		})
	}

	log.Println("Data aggregation complete, results written to output.csv")
}
