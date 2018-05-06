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
  err := vp.client.Call("PBServer.Start", argu, reply)
  log.Println(err)
  return err
}


func Make(client *rpc.Client) *VrProxy{
  vp := new(VrProxy)
  vp.client = client
  return vp
}
