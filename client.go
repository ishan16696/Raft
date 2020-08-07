package main

import (
  "fmt"
  "net"
  "os"
  "bufio"
  "log"
  "net/rpc"
  "strings"
)

type Item struct {
  Key string
  Value  string
}

var DEBUG bool=true

func main() {

  /********************************************* Simple client code ****************************************************/
  
  gateway_ip:= os.Args[1]
  gateway_portNo:= os.Args[2]
  
  gatewayAddress:= gateway_ip+":"+gateway_portNo
  
  Leader_address_portNo:=""

  //getting a connection
  connection,_:= net.Dial("tcp",gatewayAddress)

  for { 

    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')

    // send to gateway
    fmt.Fprintf(connection, text + "\n")

    message, _ := bufio.NewReader(connection).ReadString('\n')


    if DEBUG {
    	fmt.Print("Got PORT No of Leader: "+message)
    }
    

    Leader_address_portNo= ":"+message
    break

  }

  var reply Item
  var db []Item

  Leader_address1:="localhost"+Leader_address_portNo

  Leader_address:= strings.TrimSpace(string(Leader_address1))

  if DEBUG {
  	fmt.Println("Leader_address client got ",Leader_address1)
  }
  
  client, err := rpc.DialHTTP("tcp", Leader_address)

  if err != nil {
    log.Fatal("Connection error: ", err)
  }

  a := Item{"Tyagi", "199"}
  b := Item{"SAP", "1011"}
  c := Item{"Ishan", "110"}

  client.Call("API.Set", a, &reply)
  client.Call("API.Set", b, &reply)
  client.Call("API.Set", c, &reply)
  client.Call("API.GetAll", "", &db)

  fmt.Println("Database: ", db)

  client.Call("API.Edit", Item{"SAP", "102"}, &reply)
  client.Call("API.Del", c, &reply)
  client.Call("API.GetAll", "", &db)
  fmt.Println("Database: ", db)

  client.Call("API.Get", "Ishan", &reply)
  fmt.Println("Ishan : ", reply)

}