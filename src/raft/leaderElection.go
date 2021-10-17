package raft

import (
	"Raft/src/server"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// getURL returns URL
func getURL(ip string, port int, endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
}

func (r *Raft) StartElectionLoop(ctx context.Context, stopch chan struct{}) {
	// start the timer
	r.StartElectionTimer()

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing the election loop...")
			return
		case <-stopch:
			log.Println("Closing the election loop...")
			return

		case <-r.electionTimer.C:
			r.StopElectionTimer()
			r.config.Server.SetState(server.Candidate)
			log.Printf("node becomes the: %v", r.config.Server.CurrentState)
			log.Println("Triggers the election...")
			err := r.TriggerElection(stopch)
			if err != nil {
				log.Printf("Unable to become leader: %v", err)
				r.config.Server.SetState(server.Follower)
				r.StartElectionTimer()
				continue
			}
			r.SendHeartBeatToAll()

			if r.config.Server.GetState() == server.Leader {
				fmt.Println("starting the heartbeat timer ...")
				r.StartHeartBeatTimer()
			} else {
				r.StartElectionTimer()
			}
		case <-r.heartbeatTimer.C:
			r.StopHeartBeatTimer()
			err := r.SendHeartBeatToAll()
			if err != nil {
				fmt.Println(err)
			}
			r.StartHeartBeatTimer()
		}
	}
}

// TriggerElection will change the state of node and held the election
func (r *Raft) TriggerElection(stopch chan struct{}) error {
	// Increase the Term by 1.
	r.config.Server.SetTerm(r.config.Server.GetTerm() + 1)
	// candidate voted for himself
	voteReceived := 1
	// node itself is healthy
	peerHealthy := 1

	// checks the health of all nodes
	for _, port := range r.config.Server.Peers {
		healthy, err := IsServerHealthy(getURL(defaultServerIP, port, "health"))
		if !healthy || err != nil {
			log.Printf("all peers are not healthy: %v", err)
		}
		if healthy {
			peerHealthy++
		}
	}

	if peerHealthy < r.config.Server.Quorum() {
		log.Printf("No of server in healthy state: %v", peerHealthy)
		r.config.Server.SetTerm(r.config.Server.Term - 1)
		return server.QuorumLost
	}

	for _, port := range r.config.Server.Peers {
		reqVoteReply, err := r.RequestForVote(getURL(defaultServerIP, port, "askVote"))
		if err != nil {
			log.Printf("unable to get reply from server: %v", port)
		}
		fmt.Println(reqVoteReply)

		if reqVoteReply.Term > r.config.Server.GetTerm() {
			r.config.Server.SetTerm(reqVoteReply.Term)
			return server.NotLeaderError
		}
		if reqVoteReply.VoteGranted {
			voteReceived++
		}
		if voteReceived >= r.config.Server.Quorum() {
			r.config.Server.SetState(server.Leader)
			fmt.Println("node becomes the Leader")
			return nil
		}
	}
	return nil
}

func (r *Raft) SendHeartBeatToAll() error {

	for _, port := range r.config.Server.Peers {
		err := r.SendHeartBeat(getURL(defaultServerIP, port, "heartbeatPluslog"))
		if err != nil {
			log.Printf("unable to send heartbeat to server: %v", port)
		}
	}
	return nil
}

func (r *Raft) SendHeartBeat(serverURL string) error {
	var reply HeartbeatReply
	heart := Heartbeat{
		Term:        r.config.Server.GetTerm(),
		LeaderID:    r.config.Server.ServerID,
		ContainLogs: false,
	}

	dataToSend, err := json.Marshal(heart)
	if err != nil {
		log.Printf("Unable to marshal health status to json: %v", err)
		return err
	}

	response, err := http.Post(serverURL, "application/json", bytes.NewBuffer(dataToSend))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(responseData, &reply); err != nil {
		return err
	}

	fmt.Printf("HeartbeatReply: %+v\n", reply)
	return nil
}
