package raft

import(
	//"fmt"
	"sync"
	"time"
	"log"
	"math/rand"
)

type current_state string

const (

	Follower  = "follower"
	Candidate  = "candidate"
	Leader  = "Leader"
)

type Raft struct{

	mu sync.Mutex
	//cond *sync.Cond // condition variable used to avoid busy waiting 

	My_id int
	PeerIDs []int
	State current_state

	CurrentTerm int
	votedFor int

}



// type RequestVoteReply struct{
//     term int
//     follower_id int
//     voteGranted bool 
// }

var taken int=0
var done bool=false

func StartTimer(wg *sync.WaitGroup,rf *Raft){
	duration:= rand.Int()%40+200
	time.Sleep(time.Duration(duration) * time.Millisecond)
	//rf.mu.Lock()
	if taken==0{
		taken=1
		rf.TriggerElection()

	}
	//rf.mu.Unlock()
	if done==false{
		taken=0
		StartTimer(wg,rf)
	}
	wg.Done()

}





func (rf *Raft) TriggerElection(){
	rand.Seed(time.Now().UnixNano())

	rf.mu.Lock()
	rf.State = Candidate
	rf.CurrentTerm++
	rf.votedFor= rf.My_id
	term:= rf.CurrentTerm

	log.Printf("Server[%d] becomes a candidate at Term:[%d] ", rf.My_id,rf.CurrentTerm)	
	log.Printf("Server[%d] Triggers  an election at Term:[%d] ", rf.My_id,rf.CurrentTerm)
	
	votes :=1
	//done :=false
	finished :=1
	rf.mu.Unlock()
	
	for _,other_peers_id := range rf.PeerIDs{

		if other_peers_id == rf.My_id{
			continue
		}
		
		go func(id int){

			
			vote_granted := rf.CallRequestVote(id,term)

			rf.mu.Lock()
			defer rf.mu.Unlock()
			finished++

			if done{
				return
			}

			log.Printf("Server [%d] got the vote from server: [%d]",rf.My_id,id)

			if vote_granted==false{
				return
			}
			votes++

			if done==true || votes<=len(rf.PeerIDs)/2{
				return
			}else if(votes>len(rf.PeerIDs)/2){
				done=true
			}

			log.Printf("Server [%d] got the majority with votes %d in Term[%d]",rf.My_id,votes,rf.CurrentTerm)
			rf.State = Leader

		}(other_peers_id)

	}

	time.Sleep(5*time.Second)
	

}

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


// func timeOutHandler(){

// }


func (rf *Raft) CallRequestVote(server int,term int) bool{
	log.Printf("Server[%d] sending request vote to %d",rf.My_id,server)
	//time.Sleep(time.Duration(rand.Intn(100))* time.Millisecond)
	log.Printf("Server[%d] finish sending request vote to %d",rf.My_id,server)
	
	return rand.Int() %2 ==0
}


func Initialization(no_of_servers int) ([]Raft,error) {

	server := make([]Raft,no_of_servers)

	for i := 0; i < no_of_servers; i++ {
		server[i].My_id = i+1
		server[i].CurrentTerm =0
		server[i].State = Follower
		//server[i].cond = sync.NewCond(&(server[i].mu))

		for j := 0; j < no_of_servers; j++ {
			server[i].PeerIDs = append(server[i].PeerIDs,j+1)
		}
	}

	return server,nil
}
