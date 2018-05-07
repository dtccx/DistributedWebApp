package vrproxy

import (
  "common"
  "net/rpc"
  // "encoding/gob"
  "log"
)

type VrProxy struct{
  client *rpc.Client
}

func (vp *VrProxy) CallVr(argu *common.VrArgu, reply *common.VrReply) error {
  // gob.RegisterName("haha",common.SignArgs{})
  // gob.RegisterName("common.SignArgs",common.SignArgs{})
  // log.Println("reach?")


  // call primary
  err := vp.client.Call("PBServer.Start", argu, reply)
  if(err != nil) {
    // PIRMAY dies
    //If a client doesn’t receive a timely response to a request,
    //it re-sends the request to all replicas.
    //This way if the group has moved to a later view, its message will reach the new primary.
    for i := 0; i < len(common.Address); i++ {
      client_temp, err := rpc.Dial("tcp", "localhost" + common.Address[i])
      var reply *common.DealPrimayReply
      err := vp.client_temp.Call("PBServer.DealPrimay", reply)
      if(reply.OK) {
        client = client_temp
      }
    }
  }
  log.Println(err)
  return err

  client, err := rpc.Dial("tcp", "localhost:8081")

}

func (vp *VrProxy) GetPrimary() error {

}


func Make(client *rpc.Client) *VrProxy{
  vp := new(VrProxy)
  vp.client = client
  return vp
}
