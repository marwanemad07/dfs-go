package utils

import (
	"fmt"
	"strings"
	"bufio"
	"io"
	"os"
	"net"
	"log"
	"time"

	pb "dfs/proto"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/go-gota/gota/dataframe"
)

func StartTCPServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	log.Println("Data Keeper is listening for file uploads on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go HandleFileUpload(conn)
	}
}

func HandleFileUpload(conn net.Conn) {
	defer conn.Close()

	// Read filename from client
	reader := bufio.NewReader(conn)
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read filename:", err)
		return
	}
	filename = filename[:len(filename)-1] // Remove newline character

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

func StartHeartbeat() {
	for {
		SendHeartbeat("DataKeeper-1")
		time.Sleep(1 * time.Second)
	}
}

// sendHeartbeat notifies Master Tracker that this Data Keeper is alive
func SendHeartbeat(id string) {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect to Master Tracker: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.SendHeartbeat(ctx, &pb.HeartbeatRequest{DataKeeperId: id})
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
	} else {
		log.Println("Heartbeat sent successfully")
	}
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