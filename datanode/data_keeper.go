package main

import (
	"dfs/utils"
	"os"
	"fmt"
)



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <arg1> <arg2>")
		return
	}

	// Start TCP Server for receiving files
	port := os.Args[1]
	go utils.StartTCPServer(port)

	// Start gRPC heartbeat mechanism
	go utils.StartHeartbeat(port)

	select {} // Keep the process running
}