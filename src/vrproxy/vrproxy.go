package vrproxy

import (
  "common"
  "net/rpc"
)

type VrProxy struct{
  client *rpc.Client
}

func (vp *VrProxy) CallVr(argu *common.VrArgu, reply *common.VrReply) error {
  err := vp.client.Call("PBServer.Start", argu, reply)
  return err
}


func Make(client *rpc.Client) *VrProxy{
  vp := new(VrProxy)
  vp.client = client
  return vp
}
