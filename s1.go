package main

import (

    "fmt"
    "net"
    "os"
    "bufio"
    "strconv"

)

var DEBUG bool=true
var term int

// type RequestVote struct{
//     msg_type string
//     term int
//     candidate_id int
// }

type RequestVoteReply struct{
    msg_type string
    follower_id int
    voteGranted int // like 1 --- for voteGranted and 0 -- for vote not Granted
}

// type HeartBeat struct{
//   msg_type string
//   leader_id int
// }

func requestVote() string{
    text:="RequestVote"+"$"+strconv.Itoa(term)+"$"+"1"+"$"
    return text
}

func get_heartBeat() string{
  text:="HeartBeat"+"$"+"1"
  return text
}

func read(temp_message string) *RequestVoteReply {
    msg_from_server:=new(RequestVoteReply)

    var c int=0
    var count int=0

    for i:=0;i<len(temp_message);i++{

      if(string(temp_message[i])=="$"){
        
        //fmt.Println(temp_message[c:i])
        
        if count==0{
          msg_from_server.msg_type= temp_message[c:i]
          count++
        } else if count==1{
          msg_from_server.follower_id,_= strconv.Atoi(temp_message[c:i])
          count++
        }else if count==2{
          msg_from_server.voteGranted,_= strconv.Atoi(temp_message[c:i])
          count++
        }else{
          break
        }

        c=i+1

      }

    } // end of for loop


    return msg_from_server
}


func main() {

  /********************************************* Server1 code ****************************************************/
  
  fmt.Println("server1 Started......")

  server2_port:= os.Args[1]
  server3_port:= os.Args[2]

  server2Address:= "127.0.0.1"+":"+server2_port
  server3Address:= "127.0.0.1"+":"+server3_port
  

  //getting a connection
  connection2,_:= net.Dial("tcp",server2Address)
  connection3,_:= net.Dial("tcp",server3Address)
  

  term=term+1
  var voteCount int=0
  total_no_of_nodes:=3
  msg_from_server2:=new(RequestVoteReply)
  msg_from_server3:=new(RequestVoteReply)
  var isLeader bool=false
    
  fmt.Print("sending RequestVote to s2: ")
  text:= requestVote()

  // send to server2
  fmt.Fprintf(connection2, text + "\n")

  // reading a message from server2
  message2, _ := bufio.NewReader(connection2).ReadString('\n')
  temp_message:=string(message2)
  msg_from_server2= read(temp_message)
  fmt.Println(*msg_from_server2)


  fmt.Print("sending RequestVote to s3: ")
  text3:=requestVote()
  fmt.Fprintf(connection3, text3 + "\n")
    
  message3, _ := bufio.NewReader(connection3).ReadString('\n')
  temp_message=string(message3)

  msg_from_server3 = read(temp_message)
  fmt.Println(*msg_from_server3)

  if (*msg_from_server2).voteGranted==1{
    voteCount++
  }

  if (*msg_from_server3).voteGranted==1{
    voteCount++
  }

  if voteCount>total_no_of_nodes/2 {
    isLeader=true
  }

  if DEBUG{
    fmt.Println("We got the Leader...")
  }

  if isLeader==true{

    fmt.Print("sending HeartBeat to s2: ")
    text:= get_heartBeat()

    // send to server2
    fmt.Fprintf(connection2, text + "\n")

  
    fmt.Print("sending HeartBeat to s3: ")
    text3:=get_heartBeat()

    fmt.Fprintf(connection3, text3 + "\n")
    
  }


}