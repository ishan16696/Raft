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
	Follower  state = "Follower"
	Candidate state = "Candidate"
	Leader    state = "Leader"

	defaultServerIP = "127.0.0.1"
	NotVotedYet     = ""
)

// ServerConfig defines the field in server which are configurable.
type ServerConfig struct {
	ServerPort int
	TotalNodes int
	PeerPorts  []int
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
	VotedFor() string
	MemberCount() int
	QuorumSize() int
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

func (s *server) Voted() string {
	return s.VotedFor
}

func (s *server) QuorumSize() int {
	return s.Nodes/2 + 1
}

func getID(port int) string {
	portInByte := new([]byte)
	hash := sha256.New()
	hash.Write(append(*portInByte, byte(port)))
	return hex.EncodeToString(hash.Sum(*portInByte)[:20])
}

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
	}
}

func NewRaft(cfg *ServerConfig) *Raft {
	return &Raft{
		Server: NewServer(cfg),
	}
}
