package main

import(
	"fmt"
	"net"
	"bufio"
	"os"
	"strconv"
)

type RequestVote struct{
    msg_type string
    term int
    candidate_id int
}

var DEBUG bool=true

func read(temp_message string) *RequestVote{

      msg_from_server:=new(RequestVote)
      c:=0
      count:=0
      for i:=0;i<len(temp_message);i++{

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


func main() {

    fmt.Println("server3 Started......")
   
    port1:= ":"+os.Args[1]
    port2:= ":"+os.Args[2]


    listener1, _ := net.Listen("tcp",port1)
    connection1, _ := listener1.Accept()

    if DEBUG{
       fmt.Println("After1")
    }
   


    listener2, _ := net.Listen("tcp",port2) 
    connection2, _ := listener2.Accept()


    if DEBUG{
      fmt.Println("After2")    
    }
    

    
      message,_ := bufio.NewReader(connection1).ReadString('\n')
      temp_message:=string(message)

      msg_from_server1:=new(RequestVote)

      msg_from_server1=read(temp_message)

      // output message received
      fmt.Println(*msg_from_server1)
      

      // send the leader's portNo to client
      newmessage := "RequestVoteReply"+"$"+"3"+"$"+"1"+"$"

      // send new string back to client
      connection1.Write([]byte(newmessage + "\n"))


      message2, _ := bufio.NewReader(connection2).ReadString('\n')

      // output message received
      fmt.Print("Message Received from server2:", string(message2))
      

      // send new string back to client
      connection2.Write([]byte(newmessage + "\n"))

      {

        message, _ := bufio.NewReader(connection1).ReadString('\n')
        temp_message:=string(message)

        fmt.Println(temp_message)
    }
      
    

    
 

}