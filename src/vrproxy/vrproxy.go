package vrproxy

import (
  "common"
  "net/rpc"
  // "encoding/gob"
  "log"
  "strconv"
)

type VrProxy struct{
  client *rpc.Client
  addresses []string
}

func (vp *VrProxy) CallVr(argu *common.VrArgu, reply *common.VrReply) error {
  // call primary
  err := vp.client.Call("PBServer.Start", argu, reply)
  if(err != nil) {
    // PIRMAY dies
    //If a client doesn’t receive a timely response to a request,
    //it re-sends the request to all replicas.
    //This way if the group has moved to a later view, its message will reach the new primary.
    for i := 0; i < len(vp.addresses); i++ {
      client_temp, err := rpc.Dial("tcp", "localhost" + vp.addresses[i])
      if(err != nil) {
        break
      }
      var reply2 *common.DealPrimayReply
      var argu2 common.DealPrimayArgs
      client_temp.Call("PBServer.DealPrimay", argu2, reply2)
      if(reply2.OK) {
        vp = CreateVrProxy(client_temp, 8081,3)
        vp.CallVr(argu, reply)
      }
    }
  }
  log.Println(err)

  // client, err := rpc.Dial("tcp", "localhost:8081")

  return err

}



func CreateVrProxy(client *rpc.Client, startPort int, clientNumber int) *VrProxy{
  vp := new(VrProxy)
  addresses := make([]string, clientNumber)
  for i:=0; i<clientNumber; i++{
    addresses[i] = ":" + strconv.Itoa(startPort+i)
  }
  vp.client = client
  vp.addresses = addresses
  return vp
}
