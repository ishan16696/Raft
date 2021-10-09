package raft

import (
	"Raft/src/server"
	"fmt"
	"log"
)

const (
	defaultServerIP = "127.0.0.1"
)

func getURL(ip string, port int, endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
}

func (r *Raft) StartElection(stopch chan struct{}) error {
	vote := 1
	peerHealthy := 0
	r.config.Server.CurrentState = server.Candidate
	r.config.Server.Term++

	for _, port := range r.config.Server.Peers {
		healthy, err := IsServerHealthy(getURL(defaultServerIP, port, "health"))
		if !healthy || err != nil {
			log.Printf("all peers are not healthy: %v", err)
		}
		if healthy {
			peerHealthy++
		}
	}

	// if peerHealthy <= r.config.Server.Quorum()-1 {
	// 	log.Printf("No of peer server in healthy state: %v", peerHealthy)
	// 	r.config.Server.CurrentState = server.Follower
	// 	r.config.Server.Term--
	// 	return fmt.Errorf("Quorum is lost ...")
	// }

	for _, port := range r.config.Server.Peers {
		reqVoteReply, err := r.RequestForVote(getURL(defaultServerIP, port, "askVote"))
		if err != nil {
			log.Printf("unable to get reply from server: %v", port)
		}
		fmt.Println(reqVoteReply)

		fmt.Println(r.config.Server.Term)
		if reqVoteReply.Term > r.config.Server.Term {
			r.config.Server.CurrentState = server.Follower
			r.config.Server.Term = reqVoteReply.Term
			return nil
		}
		if reqVoteReply.VoteGranted {
			vote += 1
		}
		if vote >= r.config.Server.Quorum() {
			r.config.Server.CurrentState = server.Leader
			fmt.Println("I'm the leader")
			return nil
		}

	}
	return nil
}
