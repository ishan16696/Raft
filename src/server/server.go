package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

type Server interface {
	Print()
	GetServerID() string
	PeerCount() int
	Quorum() int
	ServerEndPoint() string
	GetState() string
	GetTerm() int
	GetLastVotedFor() string

	CommitIndex() uint64
	IsLogEmpty() bool
	LastCommandName() string
	ElectionTimeout() time.Duration
	SetElectionTimeout(duration time.Duration)
	HeartbeatInterval() time.Duration
	SetHeartbeatInterval(duration time.Duration)
	Init() error
	Start() error
	Stop()
	Running() bool
}

func (s *server) Print() {
	log.Printf("%+v", s)
}

func (s *server) GetServerID() string {
	return s.ServerID
}

func (s *server) Quorum() int {
	return s.Nodes/2 + 1
}

func (s *server) PeerCount() int {
	return s.Nodes - 1
}

func (s *server) ServerEndPoint() string {
	return s.EndPoint
}

func (s *server) GetState() state {
	return s.CurrentState
}

func (s *server) SetState(changedState state) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CurrentState = changedState
}

func (s *server) GetLastVotedFor() string {
	return s.VotedFor
}

func (s *server) SetLastVotedFor(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.VotedFor = id
}

func (s *server) GetTerm() int {
	return s.Term
}

func (s *server) SetTerm(newTerm int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Term = newTerm
}

// generateID takes the port no. and returns the hash of it as string.
func generateID(port int) string {
	portInByte := new([]byte)
	hash := sha256.New()
	hash.Write(append(*portInByte, byte(port)))
	return hex.EncodeToString(hash.Sum(*portInByte)[:20])
}

// getEndPoint takes input as port no. and returns the endpoint.
func getEndPoint(port int) string {
	return fmt.Sprintf("http://%s:%d", defaultServerIP, port)
}

// GetServerConfig returns the ServerConfig with default values.
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		TotalNodes: DefaultTotalNodes,
		ServerPort: DefaultServerPort,
	}
}

// NewServer returns the new server
func NewServer(cfg *ServerConfig) *server {
	return &server{
		Port:         cfg.ServerPort,
		Nodes:        cfg.TotalNodes,
		EndPoint:     getEndPoint(cfg.ServerPort),
		ServerID:     generateID(cfg.ServerPort),
		CurrentState: Follower,
		Term:         0,
		VotedFor:     NotVotedYet,
		Peers:        cfg.PeerPorts,
		mutex:        &sync.Mutex{},
	}
}

// NewRaft returns the new raft
func NewRaft(cfg *ServerConfig) *Raft {
	return &Raft{
		Server:   NewServer(cfg),
		LeaderID: NoLeader,
	}
}
