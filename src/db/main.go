package main

import(
  "log"
  "net"
  "os"
  "os/signal"
  "net/rpc"
  "strconv"
)


func createServer(clients []*rpc.Client, serverIndex int) *rpc.Client{
  ps := Make(clients, serverIndex, 0)
  server := rpc.NewServer()
  server.Register(ps)
  port := ":"+strconv.Itoa(8081+serverIndex)
  l,listenError := net.Listen("tcp", port)
  if(listenError!=nil){
    log.Println(listenError)
  }
  go server.Accept(l)
  client, err := rpc.Dial("tcp", port)
  clients[0] = client
  if(err!=nil){
    log.Println(err)
  }
  log.Println(client==nil)

  clients[serverIndex] = client

  return client
}

func main(){
  serverNum := 3
  clients := make([]*rpc.Client, serverNum)
  for i := 0; i < serverNum; i++ {
	   createServer(clients, i)
	}



  signalChan := make(chan os.Signal, 2)
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
