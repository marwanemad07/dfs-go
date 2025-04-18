package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

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

func EnsureStorageFolder(folderName string) {
	// Check if the folder exists
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		// Folder does not exist, create it with 0755 permissions
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			fmt.Printf("Failed to create folder: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Folder created: storage")
	} else {
		fmt.Printf("Folder already exists: %s\n", folderName)
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

func GetRandomIndex(length int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}

// sendResponse writes a message to the connection and flushes it
func SendResponse(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return writer.Flush()
}

// SendRequest sends multiple request fields to the server.
func SendRequest(conn net.Conn, fields ...string) error {
	if len(fields) == 0 {
		return fmt.Errorf("at least one field must be provided")
	}

	writer := bufio.NewWriter(conn)

	// Combine all fields with newline separation
	message := strings.Join(fields, "\n") + "\n"

	if _, err := writer.WriteString(message); err != nil {
		return fmt.Errorf("failed to send request (%v): %w", fields, err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush request data: %w", err)
	}

	return nil
}

func WriteFileToConnection(file *os.File, conn net.Conn) (int64, error) {
	writer := bufio.NewWriter(conn)

	// Get the file size for validation
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	expectedSize := fileInfo.Size()

	// Copy the file data to the connection
	n, err := io.Copy(writer, file)
	if err != nil {
		return n, fmt.Errorf("failed to send file data: %w", err)
	}

	// Ensure all buffered data is sent
	if err := writer.Flush(); err != nil {
		return n, fmt.Errorf("failed to flush file data: %w", err)
	}

	// Validate that the full file was sent
	if n != expectedSize {
		log.Printf("Warning: Sent %d bytes, expected %d bytes", n, expectedSize)
		return n, fmt.Errorf("incomplete data sent: sent %d bytes, expected %d", n, expectedSize)
	}

	log.Printf("Successfully sent %d bytes", n)
	return n, nil
}

func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false // Port is in use
	}
	listener.Close() // Close immediately after checking
	return true      // Port is available
}
func ShowProgress(progress <-chan int64, totalSize int64) {
	const barWidth = 50
	var received int64
	lastUpdate := time.Now()

	for r := range progress {
		received = r

		// Limit UI updates to avoid excessive printing.
		if time.Since(lastUpdate) >= 100*time.Millisecond {
			lastUpdate = time.Now()
			printProgress(received, totalSize, barWidth)
		}
	}

	// Ensure the final 100% update
	fmt.Printf("\r")
	printProgress(totalSize, totalSize, barWidth)
	fmt.Print("\n✅ Download complete!\n")
}

// printProgress displays the download progress bar and percentage.
func printProgress(received, totalSize int64, barWidth int) {
	percentage := float64(received) / float64(totalSize) * 100
	filled := int(float64(barWidth) * percentage / 100)

	bar := strings.Repeat("█", filled) + strings.Repeat("-", barWidth-filled)

	fmt.Printf("\r[%s] %.2f%% (%d/%d MB)", bar, percentage, received/1024/1024, totalSize/1024/1024)
}

func GetWiFiIPv4() (string, error) {
	// Wi-Fi interface name on Windows (common default is "Wi-Fi")
	wifiInterfaceName := "Wi-Fi" // Adjust if needed (run listInterfaces() to confirm)

	// Get the specific interface
	iface, err := net.InterfaceByName(wifiInterfaceName)
	if err != nil {
		return "", fmt.Errorf("failed to get Wi-Fi interface %s: %v", wifiInterfaceName, err)
	}

	// Get addresses for the interface
	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("failed to get addresses for %s: %v", wifiInterfaceName, err)
	}

	// Find the first IPv4 address
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			// Check if it's IPv4 and not a loopback address
			if ipNet.IP.To4() != nil && !ipNet.IP.IsLoopback() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no IPv4 address found for Wi-Fi interface %s", wifiInterfaceName)
}