package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
	"strconv"
	
)

// type RequestVoteReply struct{
//     msg_type string
//     follower_id int
//     voteGranted int // like 1 --- for voteGranted and 0 -- for vote not Granted
// }

type RequestVote struct{
    msg_type string
    term int
    candidate_id int
}

func read(temp_message string) *RequestVote{
		msg_from_server:=new(RequestVote)
		c:=0
	    count:=0
	    for i:=0;i< len(temp_message);i++{

	      if(string(temp_message[i])=="$"){
	        
	        fmt.Println(temp_message[c:i])
	        
	        if count==0{
	          msg_from_server.msg_type= temp_message[c:i]
	          count++
	        } else if count==1{
	          msg_from_server.term,_= strconv.Atoi(temp_message[c:i])
	          count++
	        }else if count==2{
	          msg_from_server.candidate_id,_= strconv.Atoi(temp_message[c:i])
	          count++
	        }else{
	          break
	        }

	        c=i+1
	      }

	    }
	    return msg_from_server
}


func main(){

	  fmt.Println("server2 Started......")
	  
	  port1:= ":"+os.Args[1]
	  server3_port:= os.Args[2]


	  server3Address:= "127.0.0.1"+":"+server3_port

	  listener1, _ := net.Listen("tcp",port1)
	  connection1, _ := listener1.Accept()
	   
	  time.Sleep(5 * time.Second)

	  //getting a connection
	  connection3,_:= net.Dial("tcp",server3Address)

	  

	  message, _ := bufio.NewReader(connection1).ReadString('\n')
	  temp_message:=string(message)

	  msg_from_server1:=new(RequestVote)

	  msg_from_server1=read(temp_message)

	  // output message received
	  fmt.Println(*msg_from_server1)

     

	  // send the leader's portNo to client
	  newmessage := "RequestVoteReply"+"$"+"2"+"$"+"1"+"$"

	  // send new string back to client
	 connection1.Write([]byte(newmessage + "\n"))
	      

	 reader := bufio.NewReader(os.Stdin)

	 fmt.Print("Text to send: ")
	  text, _ := reader.ReadString('\n')

	  // send to gateway
	  fmt.Fprintf(connection3, text + "\n")

	  message3, _ := bufio.NewReader(connection3).ReadString('\n')

	  fmt.Println("message from server3: ",message3)




	 {
	 	message, _ := bufio.NewReader(connection1).ReadString('\n')
	  	temp_message:=string(message)

	  	fmt.Println(temp_message)
	 }



	      
    
	   


}