package vrproxy

import (
  "testing"
  "fmt"
  "common"
  "net/rpc"
  "net"
  "encoding/gob"
)
func _assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
  if a == b {
    return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
  _assertEqual(t,a,b,"")
}

func assertEqualWithMsg(t *testing.T, a interface{}, b interface{}, msg string) {
  _assertEqual(t,a,b,msg)
}

func TestVrProxyCreate(t *testing.T) {
  vp := &VrProxy{}
  assertEqual(t, vp==nil, false)
}


type PBServer struct{
}

func (ts *PBServer) Start(args *common.VrArgu, reply *common.VrReply) error{
  reply.Reply = 6
	return nil
}

type GetServerNumberArgs struct{

}

type GetServerNumberReply struct{
	Number int
}

func TestCallVrTypeTransform(t *testing.T){
  gob.Register(GetServerNumberArgs{})
  gob.Register(GetServerNumberReply{})
  server := rpc.NewServer()
  ts := new(PBServer)
  server.Register(ts)
  l,listenError := net.Listen("tcp", ":9000")
  assertEqualWithMsg(t, listenError==nil, true, "t1")
	go server.Accept(l)

  client, err := rpc.Dial("tcp", ":9000")
	assertEqualWithMsg(t, err==nil, true, "t2")
	assertEqualWithMsg(t, client!=nil, true,"t3")
  vp := VrProxy{}
  vp.client = client

  innerArgu := GetServerNumberArgs{}
  argu := &common.VrArgu{}
  argu.Argu = innerArgu
  argu.Op = "Test"

  innerVrReply := GetServerNumberReply{}
	reply := &common.VrReply{}
  reply.Reply = innerVrReply
  vrErr := vp.CallVr(argu, reply)
  fmt.Println(vrErr)
  assertEqualWithMsg(t, vrErr==nil, true, "t4")

  realTypeObject, ok := reply.Reply.(int)
  assertEqualWithMsg(t, ok, true, "t5")
  assertEqualWithMsg(t, realTypeObject,6,"t6")

}


func TestCallVrTypeTransform2(t *testing.T){
  gob.Register(common.LogArgs{})
  // gob.Register(GetServerNumberReply{})
  server := rpc.NewServer()
  ts := new(PBServer)
  server.Register(ts)
  l,listenError := net.Listen("tcp", ":9001")
  assertEqualWithMsg(t, listenError==nil, true, "t1")
	go server.Accept(l)

  client, err := rpc.Dial("tcp", ":9001")
	assertEqualWithMsg(t, err==nil, true, "t2")
	assertEqualWithMsg(t, client!=nil, true,"t3")
  vp := VrProxy{}
  vp.client = client


  innerArgu := common.LogArgs{}
  argu := &common.VrArgu{}
  argu.Argu = innerArgu
  argu.Op = "Test"

  // innerVrReply := LogReply{}
  reply := &common.VrReply{}
  // reply.Reply = innerVrReply
  vrErr := vp.CallVr(argu, reply)
  fmt.Println(vrErr)
  assertEqualWithMsg(t, vrErr==nil, true, "t4")

  realTypeObject, ok := reply.Reply.(int)
  assertEqualWithMsg(t, ok, true, "t5")
  assertEqualWithMsg(t, realTypeObject,6,"t6")

}
