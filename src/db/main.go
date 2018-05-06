package main

import(
  "log"
  "net"
  "os"
  "os/signal"
  "net/rpc"
)

func main(){
  clients := make([]*rpc.Client, 1)

  peer := Make(clients, 0, 0)
  server := rpc.NewServer()
  server.Register(peer)
  l,listenError := net.Listen("tcp", ":8081")
  if(listenError!=nil){
    log.Println(listenError)
  }
  go server.Accept(l)

  client, err := rpc.Dial("tcp", ":8081")
  clients[0] = client
  if(err!=nil){
    log.Println(err)
  }
  log.Println(client==nil)

  signalChan := make(chan os.Signal, 1)
  cleanupDone := make(chan bool)
  signal.Notify(signalChan, os.Interrupt)
  go func() {
    for _ = range signalChan {
        log.Println("\nReceived an interrupt, stopping services...\n")
        cleanupDone <- true
    }
  }()
  <-cleanupDone
}
