package main

import (
  "testing"
  "net/rpc"
  "net"
  // "net/http"
  "log"
  "common"
  "strconv"
)

type Test_DB struct{
  client *rpc.Client
}

func (test_db *Test_DB) signUpUser(uname string, password string){
  args := &common.SignArgs{uname,password}
  var reply common.SignReply
  test_db.client.Call("DB.Signup", args, &reply)
}

func BuildSuiteWithPort(port int) (*DB, *rpc.Client){
  db := new(DB)
  db.user = make(map[string]User)
  db.like = make(map[string]map[int]bool)
  server := rpc.NewServer()
  server.RegisterName("DB", db)

  l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
  if e != nil {
  	log.Fatal("listen error:", e)
  }
  go server.Accept(l)

  client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
  if err != nil {
  	log.Fatal("dialing:", err)
  }

  return db,client
}


func Test_ServerSetup(t *testing.T){
  _,client := BuildSuiteWithPort(18080)
  args := &common.LogArgs{"name"}
  var reply common.LogReply
  err := client.Call("DB.Login", args, &reply)
  if err != nil {
  	log.Fatal("Test_ServerSetup: Setup server fail error:", err)
  }
  log.Println("Test_ServerSetup pass!")
}


func Test_Register(t *testing.T){
  db,client := BuildSuiteWithPort(18081)
  userName := "lala"
  passWord := "weakPW"
  args := &common.SignArgs{userName,passWord}
  var reply common.SignReply
  err := client.Call("DB.Signup", args, &reply)
  if reply.Success == false{
    log.Fatal("Test_Register rpc call fail:", err)
  }
  user := db.user[userName]
  if (user.Name != "lala") || (user.Password != "weakPW")  {
    log.Fatal("Test_Register user info wrong in db:", err)
  }
  log.Println("Test_Register pass!")
}


func Test_Login(t *testing.T){
  _,client := BuildSuiteWithPort(18082)
  userName := "lala"
  passWord := "weakPW"
  args1 := &common.SignArgs{userName,passWord}
  var reply1 common.SignReply
  client.Call("DB.Signup", args1, &reply1)

  args2 := &common.LogArgs{userName}
  var reply2 common.LogReply
  err2 := client.Call("DB.Login", args2, &reply2)
  if(reply2.Password!=passWord){
    log.Fatal("Test_Login: fail to login", err2)
  }
}




func Test_DelUser(t *testing.T){
  _,client := BuildSuiteWithPort(18083)
  userName := "lala"
  passWord := "weakPW"
  args1 := &common.SignArgs{userName,passWord}
  var reply1 common.SignReply
  client.Call("DB.Signup", args1, &reply1)

  args2 := &common.DelUserArgs{userName}
  var reply2 common.DelUserReply
  client.Call("DB.DelUser", args2, &reply2)

  args3 := &common.LogArgs{userName}
  var reply3 common.LogReply
  client.Call("DB.Login", args3, &reply3)

  if(reply3.Success==true){
    log.Fatal("Test_DelUser: deleted user still able to login")
  }
}

func (test_db *Test_DB) sendMsgFromUser(uname string, value string){
  args := &common.SendMsgArgs{uname,value}
  var reply common.SendMsgReply
  test_db.client.Call("DB.SendMsg", args, &reply)
}

func Test_SendMsg(t *testing.T){
  db,client := BuildSuiteWithPort(18084)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")

  test_db.sendMsgFromUser("u1", "value1")
  if len(db.msg)<1 {
    log.Fatal("Test_SendMsg: fail to insert message")
  }
}
