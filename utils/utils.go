package utils

import (
	"fmt"
	"strings"
	"bufio"
	"io"
	"os"
	"net"
	"log"

	"github.com/go-gota/gota/dataframe"
)

func HandleFileUpload(reader *bufio.Reader,conn net.Conn) {
	defer conn.Close()

	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read filename:", err)
		return
	}
	filename = strings.TrimSpace(filename)

	// Create the file
	file, err := os.Create("storage/" + filename)
	if err != nil {
		log.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Copy file data
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Println("Failed to save file:", err)
		return
	}

	fmt.Printf("File %s received and saved successfully!\n", filename)
}

func HandleFileDownload(reader *bufio.Reader, conn net.Conn) {
	defer conn.Close()

	// Read requested filename
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read requested filename:", err)
		return
	}
	filename = strings.TrimSpace(filename)

	// Open the file for reading
	filePath := "storage/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File not found: %s\n", filename)
		conn.Write([]byte("ERROR: File not found\n")) // Notify client
		return
	}
	defer file.Close()

	log.Println("Sending file:", filename)

	// Send confirmation to the client before sending data
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("OK\n")
	if err != nil {
		log.Println("Failed to send response:", err)
		return
	}
	writer.Flush()

	// Send file data
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Println("Error sending file:", err)
	} else {
		log.Println("File sent successfully!")
	}
}

func SendFile(address, filename string) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Connect to Data Keeper over TCP
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	// Send "UPLOAD" request type
	_, err = conn.Write([]byte("UPLOAD\n"))
	if err != nil {
		log.Fatalf("Failed to send request type: %v", err)
	}

	// Send filename first
	_, err = conn.Write([]byte(filename + "\n"))
	if err != nil {
		log.Fatalf("Failed to send filename: %v", err)
	}

	// Send file data
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatalf("Failed to send file data: %v", err)
	}

	fmt.Println("Upload complete!")
}

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

func ReceiveFile(address, filename string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connect to Data Keeper: %v", err)
		return
	}
	defer conn.Close()

	// Send "DOWNLOAD" request
	_, err = conn.Write([]byte("DOWNLOAD\n"))
	if err != nil {
		log.Fatalf("Failed to send request type: %v", err)
	}

	// Send filename
	_, err = conn.Write([]byte(filename + "\n"))
	if err != nil {
		log.Fatalf("Failed to send filename: %v", err)
	}

	// Wait for server response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if strings.HasPrefix(response, "ERROR") {
		log.Printf("Server error: %s", response)
		return
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	// Receive file data
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Printf("Error receiving file: %v", err)
		return
	}

	log.Println("File received successfully.")
}
