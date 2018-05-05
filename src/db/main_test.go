package main

import (
  "testing"
  "net/rpc"
  "net"
  // "net/http"
  "log"
  "common"
  "strconv"
  "fmt"
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
  log.Println("Test_Login pass!")
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
  log.Println("Test_DelUser pass!")
}

func (test_db *Test_DB) sendMsgFromUser(uname string, value string){
  args := &common.SendMsgArgs{uname,value}
  var reply common.SendMsgReply
  test_db.client.Call("DB.SendMsg", args, &reply)
}

func (test_db *Test_DB) getMsg() []common.Msg{
  args := &common.GetMsgArgs{}
  var reply common.GetMsgReply
  test_db.client.Call("DB.GetMsg", args, &reply)
  return reply.Msg
}

func Test_SendMsg(t *testing.T){
  db,client := BuildSuiteWithPort(18084)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")

  test_db.sendMsgFromUser("u1", "value1")
  if len(db.msg)<1 {
    log.Fatal("Test_SendMsg: fail to insert message")
  }
  log.Println("Test_SendMsg pass!")
}

func Test_GetMsg(t *testing.T){
  _,client := BuildSuiteWithPort(18085)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.sendMsgFromUser("u1", "v1")
  msgs := test_db.getMsg()
  if len(msgs)<1{
    log.Fatal("Test_GetMsg: db empty after insertion")
  }
  if(msgs[0].Value!="v1"){
    log.Fatal("Test_GetMsg: new inserted msg has wrong value")
  }
  log.Println("Test_GetMsg pass!")
}

func (test_db *Test_DB) likeMsg(uname string, msgid int){
  args := &common.LikeArgs{uname, msgid}
  var reply common.LikeReply
  test_db.client.Call("DB.LikeMsg", args, &reply)
}

func Test_LikeMsg(t *testing.T){
  db,client := BuildSuiteWithPort(18086)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.signUpUser("u2", "p1")
  test_db.sendMsgFromUser("u1", "value1")
  test_db.likeMsg("u2", 0)
  if db.like["u2"][0]!=true{
    log.Fatal("Test_LikeMsg: fail to like msg")
  }
  log.Println("Test_LikeMsg pass!")
}

func (test_db *Test_DB) isLike(uname string, msgid int) bool{
  args := &common.IsLikeArgs{uname, msgid}
  var reply common.IsLikeReply
  test_db.client.Call("DB.IsLike", args, &reply)
  return reply.Success
}

func Test_IsLike(t *testing.T){
  _,client := BuildSuiteWithPort(18087)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.signUpUser("u2", "p1")
  test_db.sendMsgFromUser("u1", "value1")
  test_db.likeMsg("u2", 0)
  b := test_db.isLike("u2", 0)
  if(b==false){
    log.Fatal("Test_IsLike: isLike fail")
  }
  log.Println("Test_IsLike pass!")
}

func Test_IsLike2(t *testing.T){
  _,client := BuildSuiteWithPort(18088)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.signUpUser("u2", "p1")
  test_db.sendMsgFromUser("u1", "value1")
  b := test_db.isLike("u2", 0)
  if(b==true){
    log.Fatal("Test_IsLike: isLike fail")
  }
  log.Println("Test_IsLike pass!")
}


func (test_db *Test_DB) unlikeMsg(uname string, msgid int) {
  args := &common.UnLikeArgs{uname, msgid}
  var reply common.UnLikeReply
  test_db.client.Call("DB.UnLikeMsg", args, &reply)
}

func Test_UnLikeMsg(t *testing.T){
  _,client := BuildSuiteWithPort(18089)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.signUpUser("u2", "p1")
  test_db.sendMsgFromUser("u1", "value1")
  test_db.likeMsg("u2", 0)
  test_db.unlikeMsg("u2", 0)
  b := test_db.isLike("u2", 0)
  if(b==true){
    log.Fatal("Test_IsLike: unlike fail")
  }
  log.Println("Test_UnLikeMsg pass!")
}

func (test_db *Test_DB) likeList(uname string) (map[int]bool, []common.Msg) {
  args := &common.LikeListArgs{uname}
  var reply common.LikeListReply
  test_db.client.Call("DB.LikeList", args, &reply)
  return reply.Lklist, reply.Msg
}

func Test_LikeList(t *testing.T){
  _,client := BuildSuiteWithPort(18090)
  test_db := Test_DB{client}
  test_db.signUpUser("u1", "p1")
  test_db.sendMsgFromUser("u1", "value1")
  test_db.sendMsgFromUser("u1", "value2")
  test_db.likeMsg("u1", 0)
  test_db.likeMsg("u1", 1)
  lkList, msgs := test_db.likeList("u1")
  assertEqual(t, len(lkList),2)
  assertEqual(t, len(msgs),2)
  assertEqual(t, true, lkList[0])
  assertEqual(t, true, lkList[1])
  log.Println("Test_LikeList pass!")
}

func createServer(i int, clients []*rpc.Client, ports []string) {
	peer := Make(clients, i, 0)
	server := rpc.NewServer()
	server.Register(peer)
	l,listenError := net.Listen("tcp", ports[i])
	if(listenError!=nil){
		log.Println(listenError)
	}
	go server.Accept(l)

	client, err := rpc.Dial("tcp", "localhost" + ports[i])
	clients[i] = client
	if(err!=nil){
		log.Println(err)
	}
	// log.Println(client==nil)

}


func Test_VrCodeSetup(t *testing.T){
  clients := make([]*rpc.Client, 3)
  srv_num := 3
  ports := []string{":8082",":8083",":8084"}

  for i := 0; i < srv_num; i++ {
  		createServer(i, clients, ports)
  }

  argu := &GetServerNumberArgs{}
	reply := &GetServerNumberReply{}
	clients[0].Call("PBServer.GetServerNumber", argu, reply)
  assertEqual(t, reply.Number, 0)
  clients[1].Call("PBServer.GetServerNumber", argu, reply)
  assertEqual(t, reply.Number, 1)
  clients[2].Call("PBServer.GetServerNumber", argu, reply)
  assertEqual(t, reply.Number, 2)

}
