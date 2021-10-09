package raft

import (
	"Raft/src/server"
	"net/http"
)

// healthCheck contains the HealthStatus of server.
type healthCheck struct {
	HealthStatus bool `json:"health"`
}

type Raft struct {
	config *server.Raft
	status int
	server *http.Server
}

type RequestVoteReply struct {
	Term int `json:"term"`
	//FollowerId  string `json:"followerID"`
	VoteGranted bool `json:"voteGranted"`
}

type RequestVote struct {
	Term        int    `json:"term"`
	CandidateId string `json:"candidateID"`
}
