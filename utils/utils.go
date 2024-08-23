package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// WriteOutput writes the aggregated data to a CSV file
func WriteOutput(data map[string]map[string]float64) {
	f, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"date", "project_id", "number_of_transactions", "total_volume_in_usd"})

	// Write data
	for date, projects := range data {
		for projectID, totalVolume := range projects {
			writer.Write([]string{
				date,
				projectID,
				"1", // Assuming each record is one transaction
				formatFloat(totalVolume),
			})
		}
	}
}

// Helper function to format float to string
func formatFloat(val float64) string {
	return fmt.Sprintf("%.2f", val)
}
