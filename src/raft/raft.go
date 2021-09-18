package raft

import (
	"Raft/src/server"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// type Raft struct {
// 	mu sync.Mutex
// 	//cond *sync.Cond // condition variable used to avoid busy waiting

// 	My_id   int
// 	PeerIDs []int
// 	State   state

// 	CurrentTerm int
// 	votedFor    int
// }

// type RequestVoteReply struct{
//     term int
//     follower_id int
//     voteGranted bool
// }

var taken int = 0
var done bool = false

// func StartTimer(wg *sync.WaitGroup, rf *Raft) {
// 	duration := rand.Int()%40 + 200
// 	time.Sleep(time.Duration(duration) * time.Millisecond)
// 	//rf.mu.Lock()
// 	if taken == 0 {
// 		taken = 1
// 		rf.TriggerElection()

// 	}
// 	//rf.mu.Unlock()
// 	if done == false {
// 		taken = 0
// 		StartTimer(wg, rf)
// 	}
// 	wg.Done()

// }

// func (rf *Raft) TriggerElection() {
// 	rand.Seed(time.Now().UnixNano())

// 	rf.mu.Lock()
// 	rf.State = Candidate
// 	rf.CurrentTerm++
// 	rf.votedFor = rf.My_id
// 	term := rf.CurrentTerm

// 	log.Printf("Server[%d] becomes a candidate at Term:[%d] ", rf.My_id, rf.CurrentTerm)
// 	log.Printf("Server[%d] Triggers  an election at Term:[%d] ", rf.My_id, rf.CurrentTerm)

// 	votes := 1
// 	//done :=false
// 	finished := 1
// 	rf.mu.Unlock()

// 	for _, other_peers_id := range rf.PeerIDs {

// 		if other_peers_id == rf.My_id {
// 			continue
// 		}

// 		go func(id int) {

// 			vote_granted := rf.CallRequestVote(id, term)

// 			rf.mu.Lock()
// 			defer rf.mu.Unlock()
// 			finished++

// 			if done {
// 				return
// 			}

// 			log.Printf("Server [%d] got the vote from server: [%d]", rf.My_id, id)

// 			if vote_granted == false {
// 				return
// 			}
// 			votes++

// 			if done == true || votes <= len(rf.PeerIDs)/2 {
// 				return
// 			} else if votes > len(rf.PeerIDs)/2 {
// 				done = true
// 			}

// 			log.Printf("Server [%d] got the majority with votes %d in Term[%d]", rf.My_id, votes, rf.CurrentTerm)
// 			rf.State = Leader

// 		}(other_peers_id)

// 	}

// 	time.Sleep(5 * time.Second)

// }

// func (rf *Raft) CallRequestVote(server int,term int) bool {

// 	log.Printf("[%d] sending request vote to %d",rf.My_id,server)

// 	args:= RequestVoteArgs{
// 		Term: term
// 		CandidateId: rf.My_id
// 	}
// 	var reply RequestVoteReply

// 	ok:= rf.sendRequestVote(rf,&args,&reply)
// 	log.Printf("[%d] finish sending request vote to %d",rf.My_id,server)

// 	if ok==false {
// 		return false
// 	}
// }

// func Initialization(no_of_servers int) ([]Raft, error) {

// 	server := make([]Raft, no_of_servers)

// 	for i := 0; i < no_of_servers; i++ {
// 		server[i].My_id = i + 1
// 		server[i].CurrentTerm = 0
// 		server[i].State = Follower
// 		//server[i].cond = sync.NewCond(&(server[i].mu))

// 		for j := 0; j < no_of_servers; j++ {
// 			server[i].PeerIDs = append(server[i].PeerIDs, j+1)
// 		}
// 	}

// 	return server, nil
// }

func (r *Raft) GetStatus() int {
	return http.StatusOK
}

func (r *Raft) serveHealth(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(r.GetStatus())
	healthCheck := &healthCheck{
		HealthStatus: func() bool {
			if r.GetStatus() == http.StatusOK {
				return true
			}
			return false
		}(),
	}
	json, err := json.Marshal(healthCheck)
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

	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", r.config.Server.Port),
		Handler: mux,
	}
	return
}

func RunRaft(ctx context.Context, cfg *server.ServerConfig) error {
	raft := new(Raft)
	raft.config = server.NewRaft(cfg)
	raft.RegisterHandler()
	if err := raft.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("Failed to start http server: %v", err)
		}
		log.Println("HTTP server closed gracefully.")
		return err
	}

	return nil
}
