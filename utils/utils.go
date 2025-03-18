package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no valid local IP found")
}

func EnsureStorageFolder() {
	folderPath := "storage"

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Folder does not exist, create it with 0755 permissions
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			fmt.Printf("Failed to create folder: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Folder created: storage")
	} else {
		fmt.Println("Folder already exists: storage")
	}
}

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}
	return dir
}

func ExtractPort(addr string) (int, error) {
	// Split the IP:Port and return the port as an integer
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	}
	return port, nil
}
