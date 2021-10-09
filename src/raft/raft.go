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
	"time"
)

func getRaftServer(cfg *server.ServerConfig) *Raft {
	return &Raft{
		config: server.NewRaft(cfg),
	}
}

// func (r *Raft) StartTimer() {
// 	duration := rand.Int()%40 + 200
// 	time.Sleep(time.Duration(duration) * time.Millisecond)
// 	//rf.mu.Lock()
// 	if taken == 0 {
// 		taken = 1

// 	}
// 	//rf.mu.Unlock()
// 	if done == false {
// 		taken = 0
// 		StartTimer(wg, rf)
// 	}
// 	wg.Done()

// }

func (r *Raft) GetStatus() int {
	return r.status
}

// IsServerHealthy checks the whether the server of given URL healthy or not.
func IsServerHealthy(serverURL string) (bool, error) {
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

func (r *Raft) RequestForVote(serverURL string) (RequestVoteReply, error) {
	var reply RequestVoteReply
	reqVote := RequestVote{
		Term:        r.config.Server.Term,
		CandidateId: r.config.Server.ServerID,
	}

	dataToSend, err := json.Marshal(reqVote)
	if err != nil {
		log.Printf("Unable to marshal health status to json: %v", err)
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

func (r *Raft) serveHealth(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(r.GetStatus())
	healthCheck := &healthCheck{
		HealthStatus: func() bool {
			return r.GetStatus() == http.StatusOK
		}(),
	}
	json, err := json.Marshal(healthCheck)
	if err != nil {
		log.Printf("Unable to marshal health status to json: %v", err)
		return
	}
	rw.Write([]byte(json))
}

func (r *Raft) sendRequestVote(rw http.ResponseWriter, req *http.Request) {
	var reqComing RequestVote
	if req.Method != "POST" {
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(reqBody, &reqComing); err != nil {
		log.Printf("Unable to unmarshal request into json: %v", err)
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
				r.config.Server.Term = reqComing.Term
				return true
			}
			return false
		}(),
	}

	json, err := json.Marshal(reqVoteReply)
	if err != nil {
		log.Printf("Unable to marshal health status to json: %v", err)
		return
	}
	rw.Write([]byte(json))
}

// RegisterHandler registers the handler for different requests
func (r *Raft) RegisterHandler() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", r.serveHealth)
	mux.HandleFunc("/askVote", r.sendRequestVote)

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
	log.Println("Registering the http request handlers...")
	r.RegisterHandler()

	log.Println("Starting the http server...")
	log.Printf("Starting HTTP server at addr: %v", r.config.Server.Port)

	r.status = http.StatusOK
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

	// let the server started..
	time.Sleep(5 * time.Second)
	log.Println("Server is started...")

	raft.StartElection(stopCh)
	<-stopCh
	return nil
}
