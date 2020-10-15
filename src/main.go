package main

import(
	"fmt"
	"Raft/src/raft"
	"log"
	//"time"
	//"math/rand"
	"sync"
)



func main(){

	var no_of_servers int
	var wg sync.WaitGroup


	fmt.Printf("Enter the no of Servers")
	fmt.Scan(&no_of_servers)

	log.Printf("Initialization of servers started")
	server,err := raft.Initialization(no_of_servers)

	if err != nil {
		log.Printf("Error occurs in initialization of Servers")
	}else{
		log.Printf("Initialization of servers completed")
	}
	

	for i := 0; i < no_of_servers; i++ {
		wg.Add(1)
		go raft.StartTimer(&wg,&server[i])
	}

	wg.Wait()

	for i := 0; i < no_of_servers; i++ {
		if server[i].State == "Leader"{
			log.Printf("Server[%d] becomes the Leader at Term[%d] ",server[i].My_id,server[i].CurrentTerm)
		}
		
	}

}