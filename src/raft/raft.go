package raft

import (
	"Raft/src/server"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// getRaftServer returns the Raft struct with all configuration.
func getRaftServer(cfg *server.ServerConfig) *Raft {
	return &Raft{
		config:         server.NewRaft(cfg),
		electionTimer:  time.NewTimer(1000 * time.Second),
		heartbeatTimer: time.NewTimer(1000 * time.Second),
		logger:         logrus.New().WithField("actor", "server"),
	}
}

func (r *Raft) GetStatus() int {
	return r.status
}

func (r *Raft) serveHealth(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(r.GetStatus())
	healthCheck := &healthCheck{
		HealthStatus: func() bool {
			return r.GetStatus() == http.StatusOK
		}(),
	}
	json, err := json.Marshal(healthCheck)
	if err != nil {
		r.logger.Errorf("Unable to marshal health status to json: %v", err)
		return
	}
	rw.Write([]byte(json))
}

func (r *Raft) sendRequestVoteReply(rw http.ResponseWriter, req *http.Request) {
	var reqComing RequestVote
	if req.Method != "POST" {
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(reqBody, &reqComing); err != nil {
		r.logger.Errorf("Unable to unmarshal request into json: %v", err)
		return
	}

	reqVoteReply := &RequestVoteReply{
		Term: func() int {
			if r.config.Server.Term < reqComing.Term {
				return reqComing.Term
			}
			return r.config.Server.Term
		}(),
		VoteGranted: func() bool {
			if r.config.Server.Term < reqComing.Term {
				r.config.Server.SetTerm(reqComing.Term)
				r.config.Server.SetState(server.Follower)
				return true
			}
			return false
		}(),
	}

	if reqVoteReply.VoteGranted {
		r.ResetElectionTimer()
	}
	json, err := json.Marshal(reqVoteReply)
	if err != nil {
		r.logger.Errorf("Unable to marshal health status to json: %v", err)
		return
	}
	rw.Write([]byte(json))
}

func (r *Raft) serveHeartbeatPluslog(rw http.ResponseWriter, req *http.Request) {
	var reqComing Heartbeat
	if req.Method != "POST" {
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(reqBody, &reqComing); err != nil {
		r.logger.Errorf("Unable to unmarshal request into json: %v", err)
		return
	}

	r.ResetElectionTimer()
	heartbeatReply := &HeartbeatReply{
		IsResetTimer: true,
		AckForLog:    false,
	}

	if reqComing.ContainLogs {
		r.logger.Errorf("take action like write ahead logs")
	}

	if r.config.Server.GetTerm() < reqComing.Term {
		r.config.Server.SetTerm(reqComing.Term)
	}

	json, err := json.Marshal(heartbeatReply)
	if err != nil {
		r.logger.Errorf("Unable to marshal health status to json: %v", err)
		return
	}
	rw.Write([]byte(json))
}

// RegisterHandler registers the handler for different requests
func (r *Raft) RegisterHandler() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", r.serveHealth)
	mux.HandleFunc("/askVote", r.sendRequestVoteReply)
	mux.HandleFunc("/heartbeatPluslog", r.serveHeartbeatPluslog)

	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", r.config.Server.Port),
		Handler: mux,
	}
}

// Stop stops the http server
func (r *Raft) Stop() error {
	return r.server.Close()
}

// StartServer start the http server
func (r *Raft) StartServer() {
	r.logger.Info("Registering the http request handlers...")
	r.RegisterHandler()

	r.logger.Info("Starting the http server...")
	r.logger.Infof("Starting HTTP server at addr: %v", r.config.Server.Port)

	if err := r.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			r.status = http.StatusInternalServerError
			log.Fatalf("Failed to start http server: %v", err)
			return
		}
		r.status = int(http.StateClosed)
		log.Println("HTTP server closed gracefully.")
	}
}

// RunRaft starts the server and timer.
func RunRaft(ctx context.Context, cfg *server.ServerConfig) error {
	stopCh := make(chan struct{})

	raft := getRaftServer(cfg)

	go raft.StartServer()
	defer raft.Stop()

	// let the server started...
	time.Sleep(3 * time.Second)
	raft.logger.Info("Server is started...")
	raft.status = http.StatusOK

	raft.StartElectionLoop(ctx, stopCh)
	return nil
}
