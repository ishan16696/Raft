package raft

import (
	"math/rand"
	"time"
)

func (r *Raft) StartElectionTimer() {
	r.electionTimer = time.NewTimer(time.Duration(rand.Intn(maximumInterval-minimumInterval)+minimumInterval) * time.Second)
}

func (r *Raft) ResetElectionTimer() {
	r.electionTimer.Reset(time.Duration(rand.Intn(maximumInterval-minimumInterval)+minimumInterval) * time.Second)
}

func (r *Raft) StopElectionTimer() {
	r.electionTimer.Stop()
}

func (r *Raft) StartHeartBeatTimer() {
	r.heartbeatTimer = time.NewTimer(DefaultHeartbeatInterval)
}

func (r *Raft) StopHeartBeatTimer() {
	r.heartbeatTimer.Stop()
}

func (r *Raft) ResetHeartBeatTimer() {
	r.electionTimer.Reset(DefaultHeartbeatInterval)
}
