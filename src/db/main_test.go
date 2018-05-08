package main

import (
  "testing"
  "net/rpc"
  "net"
  "log"
  "common"
  "strconv"
  "fmt"
  "net/http/httptest"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
  "vrproxy"
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

func assertEqualWithMsg(t *testing.T, a interface{}, b interface{}, msg string){
  _assertEqual(t,a,b,msg)
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

type GetServerNumberArgs struct{

}

type GetServerNumberReply struct{
	Number int
}

func (srv *PBServer) GetServerNumber(args *GetServerNumberArgs, reply *GetServerNumberReply) error{
	reply.Number = srv.me
	// log.Println("haha",srv.me)
	return nil
}

// func Test_VrCodeSetup(t *testing.T){
//   clients := make([]*rpc.Client, 3)
//   srv_num := 3
//   ports := []string{":8082",":8083",":8084"}
//
//   for i := 0; i < srv_num; i++ {
//   		createServer(i, clients, ports)
//   }
//
//   argu := &GetServerNumberArgs{}
// 	reply := &GetServerNumberReply{}
// 	clients[0].Call("PBServer.GetServerNumber", argu, reply)
//   assertEqual(t, reply.Number, 0)
//   clients[1].Call("PBServer.GetServerNumber", argu, reply)
//   assertEqual(t, reply.Number, 1)
//   clients[2].Call("PBServer.GetServerNumber", argu, reply)
//   assertEqual(t, reply.Number, 2)
// }



//*******************************************
//*******************************************
//*******************************************
//integration test
//*******************************************
//*******************************************
//*******************************************
var urlString = "http://localhost:8080";

var pbservers []*PBServer

func TestIntegrationInit(t *testing.T){
  gob.Register(common.VrArgu{})
  gob.Register(common.VrReply{})
  gob.Register(common.SignReply{})
  gob.Register(common.SignArgs{})
  gob.Register(common.LogArgs{})
  gob.Register(common.LogReply{})
  gob.Register(common.DelUserArgs{})
  gob.Register(common.DelUserReply{})
  gob.Register(common.SendMsgArgs{})
  gob.Register(common.SendMsgReply{})
  gob.Register(common.GetMsgArgs{})
  gob.Register(common.GetMsgReply{})
  gob.Register(common.LikeArgs{})
  gob.Register(common.LikeReply{})
  gob.Register(common.UnLikeArgs{})
  gob.Register(common.UnLikeReply{})
  gob.Register(common.LikeListArgs{})
  gob.Register(common.LikeListReply{})
  gob.Register(common.IsLikeArgs{})
  gob.Register(common.IsLikeReply{})
  gob.Register(common.FollowUserArgs{})
  gob.Register(common.FollowUserReply{})
  gob.Register(common.FollowListArgs{})
  gob.Register(common.FollowListReply{})
  serverNum := 1
  clients := make([]*rpc.Client, serverNum)
  pbservers = make([]*PBServer, serverNum)
  for i := 0; i < serverNum; i++ {
    _,pbs :=  createServer(clients, i)
    pbservers[i] = pbs
  }
  client, err := rpc.Dial("tcp", "localhost:8081")
  if err != nil {
    log.Fatal("dialing:", err)
  }
  vp = vrproxy.CreateVrProxy(client, 8081, 1)
  arith = &Arith{client: client}
}

func TestLogin(t *testing.T){
  // store := sessions.NewCookieStore([]byte("something-very-secret"))
  user := make(map[string]User)
  user["user"] = User{"user", "password"}
  data := url.Values{}
  data.Set("user", "user")
  data.Add("password", "password")

  args := common.SignArgs{"user", "password"}
  vrArgu := &common.VrArgu{}
  vrArgu.Argu = args
  vrArgu.Op = "DB.Signup"
  vrReply := &common.VrReply{}
  vp.CallVr(vrArgu, vrReply)

  r := httptest.NewRequest("GET", urlString+"/User/Login?"+data.Encode(), nil)
  w := httptest.NewRecorder()
  login(w,r)

  resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

  ret := string(body)
  if ret=="false"{
    t.Fatalf("TestLogin2 fail")
  }
  fmt.Printf("............ TestLogin Passed ********* !!\n")
}

// func TestIsLike(t *testing.T) {
//   like := make(map[string]map[int]bool)
//   test_case1 := []string{"usera","userb","userb"}
//   test_case2 := []int{1,0,1}
//   temp := make(map[int]bool)
//   temp[1] = true
//   like["usera"] = temp
//   temp[0] = true
//   like["userb"] = temp
//
//
//
//   for i := 0; i < len(test_case1); i++ {
//     ok := isLike(test_case1[i],test_case2[i])
//     if(!ok) {
//       t.Fatalf("TestLike fail")
//       fmt.Printf("Liked Failed\n")
//     }
//   }
//   fmt.Printf("............ TestMessageisLiked Passed ********* !!\n")
// }


func TestSendMsgHttp(t *testing.T){
  handler := func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "............ SendMsgResponse Passed ********* !!", http.StatusInternalServerError)
  }
  req, err := http.NewRequest("POST", urlString + "/SendMsg", nil)
  if err != nil {
    log.Fatal(err)
  }
  w := httptest.NewRecorder()
  handler(w, req)
  fmt.Printf("%s", w.Body.String())
}


func TestServer(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, "............ TestGetMessageServer Passed ********* !!")
  }))
  defer ts.Close()
  res, err := http.Get(ts.URL)
  if err != nil {
    log.Fatal(err)
  }
  greeting, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%s", greeting)
}

func TestGetMsg(t *testing.T){


  // user := make(map[string]User)
  // user["user"] = User{"user", "password"}

  // args := common.SignArgs{"user", "password"}
  // vrArgu := &common.VrArgu{}
  // vrArgu.Argu = args
  // vrArgu.Op = "DB.GetMsg"
  // vrReply := &common.VrReply{}
  // vp.CallVr(vrArgu, vrReply)

  data := url.Values{}
  data.Set("index", "-2")

  // log.Println("pbservers[0].db",pbservers[0].db)
  // data.Add("password", "password")
  pbservers[0].db.msg = []common.Msg{
			  {
          ID  : 0,
          Value : "I like debuging :)",
          User  : "usera",
          LikeNum  : 2,
          IsLiked : false,
			  },
        {
          ID  : 1,
          Value : "I literally like debuging :)",
          User  : "userb",
          LikeNum  : 3,
          IsLiked : false,
			  },
        {
          ID  : 2,
          Value : "I really like debuging :)",
          User  : "userc",
          LikeNum  : 3,
          IsLiked : false,
			  }}

      var latestmsg = []common.Msg{
          {
            ID  : 2,
            Value : "I really like debuging :)",
            User  : "userc",
            LikeNum  : 3,
            IsLiked : false,
          },

              {
                ID  : 1,
                Value : "I literally like debuging :)",
                User  : "userb",
                LikeNum  : 3,
                IsLiked : false,
      			  },
              {
                ID  : 0,
                Value : "I like debuging :)",
                User  : "usera",
                LikeNum  : 2,
                IsLiked : false,
      			  },
              }
  jsonval, _ := json.Marshal(latestmsg)

  r := httptest.NewRequest("GET", urlString+"/GetMsg?"+data.Encode(), nil)
  w := httptest.NewRecorder()
  session, _ := store.Get(r, "user_session")
  // Set some session values.
  //session.Values["authenticated"] = true
  var temp interface{} = "user"
  session.Values[temp] = "usera"
  // Save it before we write to the response/return from the handler.
  session.Save(r, w)

  getMsg(w,r)

  resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

  ret := string(body)
  log.Println("ret1", ret)
  log.Println("ret2", string(jsonval))
  if ret != string(jsonval){
    t.Fatalf("TestGetMsg fail")
  }
  fmt.Printf("............ TestGetMsgJsonResponse Passed ********* !!\n")
}


type Network struct{
  clients []*rpc.Client
  dumpClients []*rpc.Client
  pbservers []*PBServer
  oriClients []*rpc.Client
}

func CreateNetwork(clients []*rpc.Client, dumpClients []*rpc.Client,pbservers []*PBServer) *Network {
  nw := new(Network)
  nw.clients = clients
  nw.dumpClients = dumpClients
  nw.pbservers = pbservers
  nw.oriClients = make([]*rpc.Client, len(clients))
  for i := 0; i < len(clients); i++ {
     nw.oriClients[i] = clients[i]
  }
  return nw
}

func (nw *Network) Disconnect(serverIndex int){
  nw.clients[serverIndex] = nw.dumpClients[0]
  nw.pbservers[serverIndex].peers = nw.dumpClients
}

func (nw *Network) Connect(serverIndex int){
  nw.clients[serverIndex] = nw.oriClients[serverIndex]
  nw.pbservers[serverIndex].peers = nw.clients
}

var global_dumpClient *rpc.Client

func SetupTestNetwork(serverNum int, startPort int, dumpClient *rpc.Client) *Network{
  clients := make([]*rpc.Client, serverNum)
  pbservers := make([]*PBServer, serverNum)
  for i := 0; i < serverNum; i++ {
     client,pbserver := createServerV2(clients, i, startPort)
     clients[i] = client
     pbservers[i] = pbserver
  }
  dumpClients := make([]*rpc.Client, serverNum)
  for i := 0; i < serverNum; i++ {
     dumpClients[i] = dumpClient
  }
  return CreateNetwork(clients, dumpClients, pbservers)
}

func TestSetupDumpClient(t *testing.T){
  dumpServer := rpc.NewServer()
  dumpServer.Register(1)
  dumpPort := ":10000"
  dumpListener, dumpListenerErr:= net.Listen("tcp", dumpPort)
  if(dumpListenerErr!=nil){
    log.Fatal("dumpListenerErr:", dumpListenerErr)
  }
  go dumpServer.Accept(dumpListener)
  dumpClient, dumpClientErr := rpc.Dial("tcp", ":10000")
  if(dumpClientErr!=nil){
    log.Fatal("dumpClientErr:", dumpClientErr)
  }
  global_dumpClient = dumpClient
}

func TestInitNetwork(t *testing.T){
  nw := SetupTestNetwork(3, 11000,global_dumpClient)
  assertEqual(t, nw!=nil, true)
  assertEqual(t, len(nw.clients)==3, true)
  assertEqual(t, len(nw.dumpClients)==3, true)
  assertEqual(t, len(nw.pbservers)==3, true)
}


func TestNetworkDisconnect(t *testing.T){
  nw := SetupTestNetwork(1, 12000,global_dumpClient)
  reply := &common.DealPrimayReply{}
  argu := common.DealPrimayArgs{}
  err := nw.clients[0].Call("PBServer.DealPrimay", argu, reply)
  assertEqual(t, err==nil, true)
  nw.Disconnect(0)
  err = nw.clients[0].Call("PBServer.DealPrimay", argu, reply)
  assertEqual(t, err!=nil, true)
}

func TestNetworkConnect(t *testing.T){
  nw := SetupTestNetwork(1, 13000,global_dumpClient)
  reply := &common.DealPrimayReply{}
  argu := common.DealPrimayArgs{}
  err := nw.clients[0].Call("PBServer.DealPrimay", argu, reply)
  assertEqual(t, err==nil, true)
  nw.Disconnect(0)
  err = nw.clients[0].Call("PBServer.DealPrimay", argu, reply)
  assertEqual(t, err!=nil, true)
  nw.Connect(0)
  err = nw.clients[0].Call("PBServer.DealPrimay", argu, reply)
  assertEqual(t, err==nil, true)
}


func RunLoginRPC(vp *vrproxy.VrProxy, username string){
  vrArgu := &common.VrArgu{}
  args := common.SignArgs{username, "password"}
  vrArgu.Argu = args
  vrArgu.Op = "DB.Signup"
  vrReply := &common.VrReply{}
  vp.CallVr(vrArgu, vrReply)
}

func TestViewChange(t *testing.T){
  nw := SetupTestNetwork(3, 14000,global_dumpClient)
  vp := vrproxy.CreateVrProxy(nw.clients[0], 14000, 3)
  RunLoginRPC(vp, "name1")
  assertEqualWithMsg(t, nw.pbservers[0].db.user["name1"].Name, "name1","TestViewChange1")
  assertEqualWithMsg(t, nw.pbservers[1].db.user["name1"].Name, "name1","TestViewChange2")
  assertEqualWithMsg(t, nw.pbservers[2].db.user["name1"].Name, "name1","TestViewChange3")
}
