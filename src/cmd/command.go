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

const (
	DefaultServerPort int = 9090
	DefaultTotalNodes int = 3
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

// validateServerConfig validates the Server configuration.
func validateServerConfig(cfg *server.ServerConfig) error {
	if cfg.ServerPort <= 0 {
		return fmt.Errorf("server Port should be greater than 0,server port can't be: %v", cfg.ServerPort)
	}

	if cfg.TotalNodes <= 0 {
		return fmt.Errorf("total nodes should be greater than zero, total nodes can't be: %v", cfg.TotalNodes)
	}

	if cfg.TotalNodes-1 != len(cfg.PeerPorts) {
		return fmt.Errorf("no. of Peer Ports should be: %v", cfg.TotalNodes-1)
	}

	for _, port := range cfg.PeerPorts {
		if port <= 0 {
			return fmt.Errorf("peer Port should be greater than 0, peer Port can't be: %v", cfg.PeerPorts)
		}
	}

	return nil
}

//  GetServerConfig returns the server configurations.
func GetServerConfig() (*server.ServerConfig, error) {
	var peers string
	opts := new(server.ServerConfig)

	flag.IntVar(&opts.ServerPort, "port", DefaultServerPort, "specify the portNo. of node")
	flag.IntVar(&opts.TotalNodes, "nodes", DefaultTotalNodes, "total no. of nodes")
	flag.StringVar(&peers, "peers", "9091,9092", "specify the portNo. of rest of all peer nodes")
	flag.Parse()

	if err := setPeerPorts(opts, peers); err != nil {
		return opts, err
	}

	log.Printf("Server Configuration: %+v", *opts)
	return opts, nil
}

// ExecuteServer the starts the server with given configurations.
func ExecuteServer(ctx context.Context) error {
	// get the server configurations
	cfg, err := GetServerConfig()
	if err != nil {
		return err
	}

	// validation
	if err := validateServerConfig(cfg); err != nil {
		return err
	}

	// Run Server using Raft for leaderElection.
	if err := raft.RunRaft(ctx, cfg); err != nil {
		return err
	}

	return nil
}
