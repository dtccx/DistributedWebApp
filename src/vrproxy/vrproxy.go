package vrproxy

import (
  "common"
  "net/rpc"
  // "encoding/gob"
  "log"
  "strconv"
  "time"
)

type VrProxy struct{
  client *rpc.Client
  addresses []string
  view  int
}

func (vp *VrProxy) GetPrimaryAddress() string{
  // log.Println("vp.view",vp.view)
  // log.Println("len(vp.addresses)",len(vp.addresses))
  return vp.addresses[vp.view % len(vp.addresses)]
}

func (vp *VrProxy) CallVr(argu *common.VrArgu, reply *common.VrReply) error {
  // call primary
  log.Println("CallVr start")
  err := vp.client.Call("PBServer.Start", argu, reply)
  log.Println("CallVr done")
  if(err != nil) {
    log.Println("CallVr start viewchange")
    var client_new *rpc.Client
    for{
      vp.view++
      log.Println("vp increase view to:", vp.view)
      newServerAddress := vp.GetPrimaryAddress()
      client_temp, err := rpc.Dial("tcp", newServerAddress)
      if(err != nil){
        continue
      }
      log.Println("vp successfully connect to:", newServerAddress)
      vrViewChangeArgu := &common.VrViewChangeArgu{vp.view}
      vrViewChangeReply := &common.VrViewChangeReply{}
      vrViewChangeErr := client_temp.Call("PBServer.VrViewChange", vrViewChangeArgu, vrViewChangeReply)
      if(vrViewChangeErr!=nil){
        continue
      }
      log.Println("vp successfully send promptviewchange to:", newServerAddress)
      time.Sleep(2000 * time.Millisecond)
      dealPrimayArgs := common.DealPrimayArgs{}
      dealPrimayReply := &common.DealPrimayReply{}
      primaryErr := client_temp.Call("PBServer.DealPrimay", dealPrimayArgs, dealPrimayReply)
      if(primaryErr!=nil || !dealPrimayReply.OK){
        log.Println("dealPrimayReply.OK:",dealPrimayReply.OK)
        continue
      }
      log.Println("vp successfully change view to:", newServerAddress)
      client_new = client_temp
      break
    }
    vp.client = client_new
    return vp.CallVr(argu, reply)
  }
  log.Println(err)
  return err

}



func CreateVrProxy(client *rpc.Client, startPort int, clientNumber int) *VrProxy{
  vp := new(VrProxy)
  addresses := make([]string, clientNumber)
  for i:=0; i<clientNumber; i++{
    addresses[i] = "localhost:" + strconv.Itoa(startPort+i)
  }
  vp.client = client
  vp.addresses = addresses
  vp.view = 0
  return vp
}

func CreateVrProxyV2(startPort int, clientNumber int) *VrProxy{
  vp := new(VrProxy)
  addresses := make([]string, clientNumber)
  for i:=0; i<clientNumber; i++{
    addresses[i] = ":" + strconv.Itoa(startPort+i)
  }
  vp.addresses = addresses
  vp.view = 0
  client, _ := rpc.Dial("tcp", vp.GetPrimaryAddress())
  vp.client = client

  return vp
}

func MakeVrProxy(client *rpc.Client, addresses []string, view int) *VrProxy{
  vp := new(VrProxy)
  vp.client = client
  vp.addresses = addresses
  vp.view = view

  return vp
}
