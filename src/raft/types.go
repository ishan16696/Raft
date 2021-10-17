package raft

import (
	"Raft/src/server"
	"net/http"
	"time"
)

const (
	defaultServerIP = "127.0.0.1"

	maximumInterval = 8
	minimumInterval = 2

	// DefaultHeartbeatInterval is the interval that the leader will send
	// AppendEntriesRequests to followers to maintain leadership.
	DefaultHeartbeatInterval = 1 * time.Second
)

// healthCheck contains the HealthStatus of server.
type healthCheck struct {
	HealthStatus bool `json:"health"`
}

type Raft struct {
	config         *server.Raft
	status         int
	server         *http.Server
	electionTimer  *time.Timer
	heartbeatTimer *time.Timer
}

type RequestVoteReply struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"voteGranted"`
}

type RequestVote struct {
	Term        int    `json:"term"`
	CandidateId string `json:"candidateID"`
}

type Heartbeat struct {
	Term        int         `json:"term"`
	LeaderID    string      `json:"leaderID"`
	ContainLogs bool        `json:"ContainLogs"`
	Logs        interface{} `json:"logs"`
}

type HeartbeatReply struct {
	IsResetTimer bool `json:"resetTimer"`
	AckForLog    bool `json:"ack"`
}
