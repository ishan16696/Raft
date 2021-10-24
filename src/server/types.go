package server

import (
	"errors"
	"sync"
)

type state string

const (
	Follower    state = "Follower"
	Candidate   state = "Candidate"
	Leader      state = "Leader"
	UnknownSate state = "Unknown"

	NoLeader = "0"

	defaultServerIP = "127.0.0.1"
	NotVotedYet     = "NotVoted"

	DefaultServerPort int = 9090
	DefaultTotalNodes int = 3
)

var NotLeaderError = errors.New("Leader is already elected...")
var QuorumLost = errors.New("Quorum is lost...")

// ServerConfig defines the field in server which are configurable.
type ServerConfig struct {
	ServerPort int   `json:"serverPort"`
	TotalNodes int   `json:"totalNodes"`
	PeerPorts  []int `json:"peerPort"`
}

// server holds the all information for server.
type server struct {
	ServerID     string
	CurrentState state
	VotedFor     string
	Term         int
	Nodes        int
	EndPoint     string
	Port         int
	Peers        []int
	mutex        sync.Locker
}

type Raft struct {
	Server   *server
	LeaderID string
}
