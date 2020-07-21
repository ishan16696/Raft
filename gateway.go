package main

import(
	"fmt"
	"net"
	"bufio"
	"os"
	"strings"
)

type Nodes struct{
  Id string
  PortNo string
  isLeader bool
}

var List_of_Nodes []Nodes

func registerNodes(node Nodes) error{
  List_of_Nodes= append(List_of_Nodes,node)
  return nil

}

func getLeaderPortNo() string{
  var port string=""
  for _,nod:= range List_of_Nodes{
    if nod.isLeader==true{
      port=nod.PortNo
    }
  }
  return port
}

func main() {


  a1:= Nodes{"abc",os.Args[2],true}
  registerNodes(a1)

  fmt.Println("Gateway Started......")

  // listen on all interfaces
  portNo:= ":"+os.Args[1]

  listener, _ := net.Listen("tcp", portNo)

  conn, _ := listener.Accept()

  for {

    message, error := bufio.NewReader(conn).ReadString('\n')

    if error!=nil{
    	fmt.Println(error)
    	return
    }

    // output message received
    fmt.Print("Message Received:", string(message))
    if strings.TrimSpace(string(message)) == "STOP" {
    	break
    }

    // send the leader's portNo to client
    newmessage := getLeaderPortNo() 

    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
    break

  }
 

}