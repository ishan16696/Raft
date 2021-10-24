package raft

import (
	"Raft/src/server"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// getURL returns URL
func getURL(ip string, port int, endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
}

func (r *Raft) StartElectionLoop(ctx context.Context, stopch chan struct{}) {
	// start the timer
	r.StartElectionTimer()
	log := logrus.New().WithField("actor", "leader-elector")

	for {
		select {
		case <-ctx.Done():
			log.Info("Closing the election loop...")
			return
		case <-stopch:
			log.Info("Closing the election loop...")
			return

		case <-r.electionTimer.C:
			// stop the election timer and become "Candidate".
			// And trigger the election to become "Leader".

			r.StopElectionTimer()
			r.config.Server.SetState(server.Candidate)
			log.Infof("node becomes the: %v", r.config.Server.GetState())
			log.Info("Triggers the election...")
			err := r.TriggerElection()
			if err != nil {
				log.Errorf("Unable to become leader: %v", err)
				r.config.Server.SetState(server.Follower)
				r.StartElectionTimer()
				continue
			}
			r.SendHeartBeatToAll()

			if r.config.Server.GetState() == server.Leader {
				log.Info("starting the heartbeat timer ...")
				r.StartHeartBeatTimer()
			} else {
				r.config.Server.SetState(server.Follower)
				r.StartElectionTimer()
			}

		case <-r.heartbeatTimer.C:
			r.StopHeartBeatTimer()
			err := r.SendHeartBeatToAll()
			if err != nil {
				log.Errorf("unable to heartbeat: %v", err)
			}
			r.StartHeartBeatTimer()
		}
	}
}

// TriggerElection will change the state of node and held the election
func (r *Raft) TriggerElection() error {
	// Increase the Term by 1.
	r.config.Server.SetTerm(r.config.Server.GetTerm() + 1)
	// candidate voted for itself.
	voteReceived := 1
	// node itself is healthy.
	peerHealthy := 1

	// checks the health of all nodes
	for _, port := range r.config.Server.Peers {
		healthy, err := isServerHealthy(getURL(defaultServerIP, port, "health"))
		if !healthy || err != nil {
			r.logger.Errorf("all peers are not healthy: %v", err)
		}
		if healthy {
			peerHealthy++
		}
	}

	if peerHealthy < r.config.Server.Quorum() {
		r.logger.Infof("No of server in healthy state: %v", peerHealthy)
		r.config.Server.SetTerm(r.config.Server.Term - 1)
		return server.QuorumLost
	}

	for _, port := range r.config.Server.Peers {
		reqVoteReply, err := r.requestForVote(getURL(defaultServerIP, port, "askVote"))
		if err != nil {
			r.logger.Errorf("unable to get reply from server: %v", port)
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
			r.logger.Info("node becomes the Leader")
			return nil
		}
	}
	return server.NotLeaderError
}

// isServerHealthy checks the whether the server of given URL healthy or not.
func isServerHealthy(serverURL string) (bool, error) {
	var health healthCheck

	response, err := http.Get(serverURL)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(responseData, &health); err != nil {
		return false, err
	}
	return health.HealthStatus, nil
}

// requestForVote sends a "POST" request for asking vote.
func (r *Raft) requestForVote(serverURL string) (RequestVoteReply, error) {
	var reply RequestVoteReply
	reqVote := RequestVote{
		Term:        r.config.Server.Term,
		CandidateId: r.config.Server.ServerID,
	}

	dataToSend, err := json.Marshal(reqVote)
	if err != nil {
		r.logger.Errorf("Unable to marshal health status to json: %v", err)
		return reply, err
	}

	response, err := http.Post(serverURL, "application/json", bytes.NewBuffer(dataToSend))
	if err != nil {
		return reply, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return reply, err
	}

	if err := json.Unmarshal(responseData, &reply); err != nil {
		return reply, err
	}
	return reply, nil
}

func (r *Raft) SendHeartBeatToAll() error {

	for _, port := range r.config.Server.Peers {
		err := r.SendHeartBeat(getURL(defaultServerIP, port, "heartbeatPluslog"))
		if err != nil {
			r.logger.Errorf("unable to send heartbeat to server: %v", port)
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
		r.logger.Errorf("Unable to marshal health status to json: %v", err)
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

	r.logger.Infof("HeartbeatReply: %+v\n", reply)
	return nil
}
