package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"        
	"log"
	"os"
	"time"
	"strconv"
	"strings"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"crypto-aggregator/coingecko"
	"crypto-aggregator/utils"
)

type AggregatedData struct {
	ProjectID            string
	Date                 string
	NumberOfTransactions int
	TotalVolumeUSD       float64
}

type BigQueryRow struct {
	Date                 string  `bigquery:"date"`
	ProjectID            string  `bigquery:"project_id"`
	NumberOfTransactions int     `bigquery:"number_of_transactions"`
	TotalVolumeUSD       float64 `bigquery:"total_volume_usd"`
}

func main() {
	// Get environment variables for configuration
	bucketName := os.Getenv("BUCKET_NAME")
	projectID := os.Getenv("PROJECT_ID")
	bqDataset := os.Getenv("BQ_DATASET")
	bqTable := os.Getenv("BQ_TABLE")
	
	processDate := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")


	// Initialize context
	ctx := context.Background()

	// Initialize GCS client
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer storageClient.Close()

	// Open the input file from GCS
	rc, err := storageClient.Bucket(bucketName).Object(fmt.Sprintf("%s/input.csv", processDate)).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to open GCS object: %v", err)
	}
	defer rc.Close()

	// Create a CSV reader
	reader := csv.NewReader(rc)
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
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read record: %v", err)
		}

		ts := record[1]        // Assuming this is the column for timestamp
		projectID := record[3] // Assuming this is the column for project ID
		propsStr := record[14] // Assuming this is the column for props JSON
		numsStr := record[15]  // Assuming this is the column for nums JSON

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
				ProjectID:            projectID,
				Date:                 date,
				NumberOfTransactions: 1,
				TotalVolumeUSD:       volumeUSD,
			}
		}
	}

	// Initialize BigQuery client
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer bqClient.Close()

	// Write the aggregated data to BigQuery
	u := bqClient.Dataset(bqDataset).Table(bqTable).Uploader()

	var rows []*BigQueryRow
	for _, data := range aggregatedData {
		rows = append(rows, &BigQueryRow{
			Date:                 data.Date,
			ProjectID:            data.ProjectID,
			NumberOfTransactions: data.NumberOfTransactions,
			TotalVolumeUSD:       data.TotalVolumeUSD,
		})
	}

	if err := u.Put(ctx, rows); err != nil {
		log.Fatalf("Failed to insert data into BigQuery: %v", err)
	}

	log.Println("Data aggregation complete, results written to BigQuery")
}
