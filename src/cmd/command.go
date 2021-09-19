package cmd

import (
	"Raft/src/raft"
	"Raft/src/server"
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func setPeerPorts(cfg *server.ServerConfig, peers string) error {
	peerPorts := strings.Split(peers, ",")

	for i := 0; i < len(peerPorts); i++ {
		temp, err := strconv.Atoi(peerPorts[i])
		if err != nil {
			return err
		}
		cfg.PeerPorts = append(cfg.PeerPorts, temp)
	}
	return nil
}

func GetServerConfig() (*server.ServerConfig, error) {
	var peers string
	opts := new(server.ServerConfig)

	flag.IntVar(&opts.ServerPort, "port", 9090, "specify the portNo.")
	flag.IntVar(&opts.TotalNodes, "nodes", 3, "total no. of nodes")
	flag.StringVar(&peers, "peers", "9091,9092", "specify the portNo. of all peers")
	flag.Parse()

	if err := setPeerPorts(opts, peers); err != nil {
		return opts, err
	}
	log.Println(opts)
	return opts, nil
}

func validateServerConfig(cfg *server.ServerConfig) error {
	if cfg.ServerPort <= 0 {
		return fmt.Errorf("Server Port should be greater than 0.")
	}

	if cfg.TotalNodes <= 0 {
		return fmt.Errorf("Total nodes should be greater than 0")
	}

	if cfg.TotalNodes-1 != len(cfg.PeerPorts) {
		return fmt.Errorf("No. of Peer Ports should be: %v", cfg.TotalNodes-1)
	}

	for _, port := range cfg.PeerPorts {
		if port <= 0 {
			return fmt.Errorf("Peer Port should be greater than 0.")
		}
	}

	return nil
}

func ExecuteServer(ctx context.Context) error {
	// get the server configurations
	cfg, err := GetServerConfig()
	if err != nil {
		return err
	}

	if err := validateServerConfig(cfg); err != nil {
		return err
	}

	// Run Server using Raft for leaderElection.
	if err := raft.RunRaft(ctx, cfg); err != nil {
		return err
	}

	return nil
}
