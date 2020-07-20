package main

import(
	"fmt"
	"net"
	//"bufio"
	"os"
	//"strings"
  //"net"
  "net/http"
  "net/rpc"
  "log"
)

type Item struct {
  Key string
  Value  string
}

type API int

var database []Item

func (a *API) GetAll(empty string, reply *[]Item) error {
  *reply = database
  return nil
}

func (a *API) Get(title string, reply *Item) error {
  var getItem Item

  for _, val := range database {
    if val.Key == title {
      getItem = val
    }
  }

  *reply = getItem

  return nil
}

func (a *API) Set(item Item, reply *Item) error {
  database = append(database, item)
  *reply = item
  return nil
}

func (a *API) Edit(item Item, reply *Item) error {
  var changed Item

  for idx, val := range database {
    if val.Key == item.Key {
      database[idx] = Item{item.Key, item.Value}
      changed = database[idx]
    }
  }

  *reply = changed
  return nil
}

func (a *API) Del(item Item, reply *Item) error {
  var del Item

  for idx, val := range database {
    if val.Key == item.Key && val.Value == item.Value {
      database = append(database[:idx], database[idx+1:]...)
      del = item
      break
    }
  }

  *reply = del
  return nil
}

func main() {


  fmt.Println("Leader Started......")

  // listen on all interfaces
  portNo:= ":"+os.Args[1]
  api := new(API)
  err := rpc.Register(api)
  if err != nil {
    log.Fatal("error registering API", err)
  }

  rpc.HandleHTTP()

  listener, err := net.Listen("tcp", portNo)

  if err != nil {
    log.Fatal("Listener error", err)
  }
  log.Printf("serving rpc on port %v ", portNo)

  http.Serve(listener, nil)

  if err != nil {
    log.Fatal("error serving: ", err)
  }

  //conn.close()

}