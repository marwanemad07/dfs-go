package main

import (
	"dfs/utils"
)



func main() {
	// Start TCP Server for receiving files
	go utils.StartTCPServer(":50051")

	// Start gRPC heartbeat mechanism
	go utils.StartHeartbeat()

	select {} // Keep the process running
}