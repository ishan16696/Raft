package cmd

import (
	"Raft/src/raft"
	"Raft/src/server"
	"context"
	"flag"
	"log"
)

func setPeerPorts(opts *server.ServerConfig) {
	for i := 1; i < opts.TotalNodes; i++ {
		opts.PeerPorts = append(opts.PeerPorts, opts.ServerPort+uint(i))
	}
}

func GetServerConfig() *server.ServerConfig {
	opts := new(server.ServerConfig)

	flag.UintVar(&opts.ServerPort, "port", 9090, "specify the portNo.")
	flag.IntVar(&opts.TotalNodes, "nodes", 3, "total no. of nodes")
	flag.Parse()

	setPeerPorts(opts)
	log.Println(opts)
	return opts
}

func ExecuteServer(ctx context.Context) error {

	cfg := GetServerConfig()

	if err := raft.RunRaft(ctx, cfg); err != nil {
		return err
	}

	return nil
}
