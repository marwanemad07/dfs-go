package utils

import (
	"fmt"
	"strings"
	"github.com/go-gota/gota/dataframe"
)

func PrintDataFrame(df dataframe.DataFrame) {
	records := df.Records() // Get all rows as [][]string

	if len(records) == 0 {
		fmt.Println("Empty DataFrame")
		return
	}

	// Print header
	header := records[0]
	fmt.Printf("| %-40s | %-20s | %-30s | %-10s |\n", header[0], header[1], header[2], header[3])
	fmt.Println(strings.Repeat("-", 110))

	// Print rows
	for _, row := range records[1:] { // Skip header row
		fmt.Printf("| %-40s | %-20s | %-30s | %-10s |\n", row[0], row[1], row[2], row[3])
	}
	fmt.Println(strings.Repeat("-", 110))
}