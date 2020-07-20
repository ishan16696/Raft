package main

import(
	"fmt"
	"net"
	"bufio"
	"os"
	"strings"
)


func main() {


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
    newmessage := os.Args[2]

    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
    break

  }
 

}