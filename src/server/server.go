package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

type state string

const (
	Follower    state = "Follower"
	Candidate   state = "Candidate"
	Leader      state = "Leader"
	UnknownSate state = "Unknown"

	defaultServerIP = "127.0.0.1"
	NotVotedYet     = "NotVoted"
)

// ServerConfig defines the field in server which are configurable by user.
type ServerConfig struct {
	ServerPort int   `json:"serverPort"`
	TotalNodes int   `json:"totalNodes"`
	PeerPorts  []int `json:"peerPort"`
}

type server struct {
	ServerID     string
	CurrentState state
	VotedFor     string
	Term         int
	Nodes        int
	EndPoint     string
	Port         int
	Peers        []int
}

type Raft struct {
	Server *server
	mutex  sync.Locker
}

type Server interface {
	GetServerID() string
	LeaderID() string

	Term() uint64
	CommitIndex() uint64
	LastVotedFor() string
	MemberCount() int
	Quorum() int
	IsLogEmpty() bool
	LastCommandName() string
	GetState() string
	ElectionTimeout() time.Duration
	SetElectionTimeout(duration time.Duration)
	HeartbeatInterval() time.Duration
	SetHeartbeatInterval(duration time.Duration)
	Init() error
	Start() error
	Stop()
	Running() bool
	Print()
}

func (s *server) Print() {
	log.Printf("%+v", s.ServerID)
}

func (s *server) GetState() state {
	return s.CurrentState
}

func (s *server) GetServerID() string {
	return s.ServerID
}

func (s *server) LastVotedFor() string {
	return s.VotedFor
}

func (s *server) Quorum() int {
	return s.Nodes/2 + 1
}

// getID takes the port no. and returns the hash of it as string.
func getID(port int) string {
	portInByte := new([]byte)
	hash := sha256.New()
	hash.Write(append(*portInByte, byte(port)))
	return hex.EncodeToString(hash.Sum(*portInByte)[:20])
}

// getEndPoint takes input as port no. and returns the endpoint.
func getEndPoint(port int) string {
	return fmt.Sprintf("http://%s:%d", defaultServerIP, port)
}

// NewServer returns the new server
func NewServer(cfg *ServerConfig) *server {
	return &server{
		Port:         cfg.ServerPort,
		Nodes:        cfg.TotalNodes,
		EndPoint:     getEndPoint(cfg.ServerPort),
		ServerID:     getID(cfg.ServerPort),
		CurrentState: Follower,
		Term:         0,
		VotedFor:     NotVotedYet,
		Peers:        cfg.PeerPorts,
	}
}

func NewRaft(cfg *ServerConfig) *Raft {
	return &Raft{
		Server: NewServer(cfg),
	}
}
