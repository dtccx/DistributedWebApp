package main

import (
  "testing"
  "net/rpc"
  "net"
  "net/http"
  "log"
  "common"
  "strconv"
)


func BuildSuiteWithPort(port int) (*DB, *rpc.Client){
  db := new(DB)
  db.user = make(map[string]User)
  db.like = make(map[string]map[int]bool)
  rpc.Register(db)
  rpc.HandleHTTP()

  l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
  if e != nil {
  	log.Fatal("listen error:", e)
  }
  go http.Serve(l, nil)

  client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(port))
  if err != nil {
  	log.Fatal("dialing:", err)
  }

  return db,client
}


func Test_ServerSetup(t *testing.T){
  _,client := BuildSuiteWithPort(8080)

  args := &common.LogArgs{"name"}
  var reply common.LogReply
  err := client.Call("DB.Login", args, &reply)
  if err != nil {
  	log.Fatal("arith error:", err)
  }
  log.Println("Test_ServerSetup pass!")
}
