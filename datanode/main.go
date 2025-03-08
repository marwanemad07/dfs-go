package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "dfs/proto" // update this import path
)

// dataKeeperServer implements Data Keeper logic.
type dataKeeperServer struct {
	id string
}

// NewDataKeeperServer returns a new Data Keeper server.
func NewDataKeeperServer(id string) *dataKeeperServer {
	return &dataKeeperServer{id: id}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	// You might register DataKeeper-specific RPCs here if needed.
	// For now, the Data Keeper will focus on sending heartbeats.
	// pb.RegisterDataKeeperServer(grpcServer, NewDataKeeperServer("DataKeeper-1"))

	log.Println("Data Keeper is listening on port 50051")

	// Start sending heartbeats in a separate goroutine.
	go func() {
		for {
			sendHeartbeat("DataKeeper-1")
			time.Sleep(1 * time.Second)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// sendHeartbeat connects to the Master Tracker and sends a heartbeat.
func sendHeartbeat(id string) {
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
		log.Printf("Heartbeat error from %s: %v", id, err)
	} else {
		log.Printf("Heartbeat sent from %s", id)
	}
}
