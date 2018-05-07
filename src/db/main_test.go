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

  //for test
  "labrpc"
  "sync"
	"time"
  "errors"
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

type GetServerNumberArgs struct{

}

type GetServerNumberReply struct{
	Number int
}

func (srv *PBServer_test) GetServerNumber(args *GetServerNumberArgs, reply *GetServerNumberReply) error{
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
  gob.Register(common.TestReply{})
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
  vp = vrproxy.Make(client)
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

//_______________________________________________________________________
//______________________________test for replication
//_______________________________________________________________________
func MakeTest(peers []*labrpc.ClientEnd, me int, startingView int) *PBServer_test {
	gob.Register(common.TestReply{})
	srv := &PBServer_test{
		peers:     	 		peers,
		me:             me,
		currentView:    startingView,
		lastNormalView: startingView,
		status:         NORMAL,
	}
	// all servers' log are initialized with a dummy command at index 0
	var v interface{}
	srv.log = append(srv.log, v)

	// Your other initialization code here, if there's any
	return srv
}

func (srv *PBServer_test) IsCommitted(index int) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.commitIndex >= index {
		return true
	}
	return false
}

func (srv *PBServer_test) ViewStatus() (currentView int, statusIsNormal bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView, srv.status == NORMAL
}

func (srv *PBServer_test) GetEntryAtIndex(index int) (ok bool, command interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.log) > index {
		return true, srv.log[index]
	}
	return false, command
}

func (srv *PBServer_test) Kill() {
	// Your code here, if necessary
}


func (srv *PBServer_test) DealPrimay(args common.DealPrimayArgs, reply *common.DealPrimayReply) error {
	if (srv.me != GetPrimary(srv.currentView, len(srv.peers))){
		reply.OK = false
	} else {
		reply.OK = true
	}
	return nil
}

func (srv *PBServer_test) Start(args common.VrArgu, reply *common.VrReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// do not process command if status is not NORMAL
	// and if i am not the primary in the current view
	test000 := GetPrimary(srv.currentView, len(srv.peers))
	log.Print("GetPrimary(srv.currentView, len(srv.peers))",test000, srv.me)
	if srv.status != NORMAL {
		if(args.Op == "test"){
			reply2 := common.TestReply{-1, srv.currentView, false}
			// log.Print("LikeList:", reply2)
			reply.Reply = reply2
		}
		return errors.New("status is INNORMAL")
	}	else if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
		if(args.Op == "test"){
			// temp, _ := args.Argu.(common.LikeListArgs)
			// log.Print(temp)
			// var reply2 common.LikeListReply
			// srv.db.LikeList(&temp, &reply2)

			reply2 := common.TestReply{-1, srv.currentView, false}
			// log.Print("LikeList:", reply2)
			reply.Reply = reply2
		}
		return errors.New("This is not Primary SRV")
	}

	command := args
	//append the command in its log
	srv.log = append(srv.log, command)
	srv.resendPrepare(command, len(srv.log) - 1, srv.currentView, srv.commitIndex)

	//write -----> use
	op := args.Op
	switch op{
	case "DB.Login":
		//var temp *common.LogArgs
		temp,  _ := args.Argu.(common.LogArgs)
		var reply2 common.LogReply
		srv.db.Login(&temp, &reply2)
		reply.Reply = reply2
	case "DB.Signup":
		temp, _ := args.Argu.(common.SignArgs)
		log.Print(temp)
		var reply2 common.SignReply
		srv.db.Signup(&temp, &reply2)
		log.Print("sign")
		reply.Reply = reply2
	case "DB.DelUser":
		temp, _ := args.Argu.(common.DelUserArgs)
		log.Print(temp)
		var reply2 common.DelUserReply
		srv.db.DelUser(&temp, &reply2)
		log.Print("delete user")
		reply.Reply = reply2
	case "DB.SendMsg":
		temp, _ := args.Argu.(common.SendMsgArgs)
		log.Print(temp)
		var reply2 common.SendMsgReply
		srv.db.SendMsg(&temp, &reply2)
		log.Print("snedmsg")
		reply.Reply = reply2
	case "DB.GetMsg":
		temp, _ := args.Argu.(common.GetMsgArgs)
		log.Print(temp)
		var reply2 common.GetMsgReply
		srv.db.GetMsg(&temp, &reply2)
		log.Print("getmsg")
		reply.Reply = reply2
	case "DB.LikeMsg":
		temp, _ := args.Argu.(common.LikeArgs)
		log.Print(temp)
		var reply2 common.LikeReply
		srv.db.LikeMsg(&temp, &reply2)
		log.Print("like")
		reply.Reply = reply2
	case "DB.UnLikeMsg":
		temp, _ := args.Argu.(common.UnLikeArgs)
		log.Print(temp)
		var reply2 common.UnLikeReply
		srv.db.UnLikeMsg(&temp, &reply2)
		log.Print("unlike")
		reply.Reply = reply2
	case "DB.IsLike":
		temp, _ := args.Argu.(common.IsLikeArgs)
		log.Print(temp)
		var reply2 common.IsLikeReply
		srv.db.IsLike(&temp, &reply2)
		log.Print("IsLike")
		reply.Reply = reply2
	case "DB.FollowList":
		temp, _ := args.Argu.(common.FollowListArgs)
		log.Print(temp)
		var reply2 common.FollowListReply
		srv.db.FollowList(&temp, &reply2)
		log.Print("FollowList")
		reply.Reply = reply2
	case "DB.FollowUser":
		temp, _ := args.Argu.(common.FollowUserArgs)
		log.Print(temp)
		var reply2 common.FollowUserReply
		srv.db.FollowUser(&temp, &reply2)
		log.Print("FollowUser")
		reply.Reply = reply2
	case "DB.LikeList":
		temp, _ := args.Argu.(common.LikeListArgs)
		log.Print(temp)
		var reply2 common.LikeListReply
		srv.db.LikeList(&temp, &reply2)
		log.Print("LikeList:", reply2)
		reply.Reply = reply2
	case "test":
		reply2 := common.TestReply{len(srv.log) - 1, srv.currentView, true}
		// log.Print("LikeList:", reply2)
		reply.Reply = reply2
		//reply.Reply = {len(srv.log) - 1, srv.currentView, true}
	}



	log.Println("return")
//	log.Println(srv.IsCommitted(index))
	return nil
}

func (srv *PBServer_test) resendPrepare(command interface{}, index int, currentView int, commitIndex int) {
	log.Println("start resendPrepare")
	replys := make([]*PrepareReply, len(srv.peers))
	for i, _ := range srv.peers {
		if i == srv.me {
			// replys[i] = &PrepareReply{
			// 	View:				srv.currentView,
			// 	Success:		true,
			// }
			continue
		}
		prepareArgs := &PrepareArgs{
			View:         currentView,         // the primary's current view
			PrimaryCommit:commitIndex, 			   // the primary's commitIndex
			Index:        index,        // the index position at which the log entry is to be replicated on backups
			Entry:        command,								 // the log entry to be replicated
		}
		replys[i] = &PrepareReply{}
		srv.sendPrepare(i, prepareArgs, replys[i])

	}
	//Recovery

	//change committedindex
	suc_num := 0
	for i, reply := range replys {
		if i == srv.me {
			continue
		}
		//log.Println(reply)
		if reply.Success == true {
			suc_num += 1
		}
	}
	//log.Printf("[%d] successful number", suc_num)
	log.Println("suc_num:", suc_num)
	log.Println("len(srv.peers)", len(srv.peers))
	if suc_num >= (len(srv.peers) - 1) / 2 {
		for {
			// srv.mu.Lock()
			if srv.commitIndex + 1 == index {
				srv.commitIndex += 1
				// srv.mu.Unlock()
				break
			}
			// srv.mu.Unlock()
		}
	} else {
			//if not committed, resend prepare until index commmit
			log.Println("resendPrepare call itself")
			go srv.resendPrepare(command, index, currentView, commitIndex)
	}

	log.Println("end resendPrepare")
}

func (srv *PBServer_test) sendPrepare(server int, args *PrepareArgs, reply *PrepareReply) bool {
	log.Println("sendPrepare:", srv.peers[server])
	log.Println("sendPrepare num:", len(srv.peers))
	prepareErr := srv.peers[server].Call("PBServer_test.Prepare", args, reply)
	log.Println("prepareErr:",prepareErr)
	return prepareErr
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer_test) Prepare(args *PrepareArgs, reply *PrepareReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()

	//whether the next entry to be added to the log is indeed at the index specified in the message
	//srv.log[args.Index] == args.Entry??

	//log.Printf("[%d] Server get prepare call", srv.me)
	//log.Printf("[%d] currentView [%d] prepareview", srv.currentView, args.View)

	if srv.currentView == args.View && len(srv.log) == args.Index {
		srv.log = append(srv.log, args.Entry)
		reply.View = srv.currentView
		reply.Success = true
		srv.commitIndex = args.PrimaryCommit
		srv.lastNormalView = srv.currentView
		return
	}	else if srv.currentView == args.View && len(srv.log) > args.Index {
		//log's capaticy is bigger, so should return true, but not append,and not recover
		reply.View = args.View
		reply.Success = true
		return
	} else {
		//reply.View = srv.currentView
		reply.Success = false
		//no way to recover
		if(srv.currentView > args.View) {
			return
		}

		//deal with other prepare, no improve, don't know why.
		//go srv.sendPrepare(srv.me, args, reply)

		srv.status = RECOVERING
		go func() {
			recoveryArgs := &RecoveryArgs {
				View:			args.View,
				Server:		srv.me,
			}
			recoveryReply := &RecoveryReply{}
			ok := srv.peers[GetPrimary(args.View, len(srv.peers))].Call("PBServer.Recovery", recoveryArgs, recoveryReply)
			if ok == true && recoveryReply.Success == true {
				srv.currentView = recoveryReply.View
				srv.lastNormalView = recoveryReply.View
				srv.commitIndex = recoveryReply.PrimaryCommit
				srv.log = recoveryReply.Entries
				srv.status = NORMAL
			}
		}()
	}
	// outdate := srv.currentView < args.View
	// log.Printf("if the date outdated",outdate)
	// if outdate {
	//
	// }


}

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer_test) Recovery(args *RecoveryArgs, reply *RecoveryReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.status == NORMAL {
		reply.View = srv.currentView
		reply.Entries = srv.log
		reply.PrimaryCommit = srv.commitIndex
		reply.Success = true
	} else {
		reply.Success = false
	}

}


func (srv *PBServer_test) PromptViewChange(newView int) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	newPrimary := GetPrimary(newView, len(srv.peers))

	if newPrimary != srv.me { //only primary of newView should do view change
		return
	} else if newView <= srv.currentView {
		return
	}
	vcArgs := &ViewChangeArgs{
		View: newView,
	}
	vcReplyChan := make(chan *ViewChangeReply, len(srv.peers))
	// send ViewChange to all servers including myself
	for i := 0; i < len(srv.peers); i++ {
		go func(server int) {
			var reply ViewChangeReply
			ok := srv.peers[server].Call("PBServer_test.ViewChange", vcArgs, &reply)
			// fmt.Printf("node-%d (nReplies %d) received reply ok=%v reply=%v\n", srv.me, nReplies, ok, r.reply)
			if ok {
				vcReplyChan <- &reply
			} else {
				vcReplyChan <- nil
			}
		}(i)
	}

	// wait to receive ViewChange replies
	// if view change succeeds, send StartView RPC
	go func() {
		var successReplies []*ViewChangeReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil && r.Success {
				successReplies = append(successReplies, r)
			}
			if nReplies == len(srv.peers) || len(successReplies) == majority {
				break
			}
		}
		ok, log := srv.determineNewViewLog(successReplies)
		if !ok {
			return
		}
		svArgs := &StartViewArgs{
			View: vcArgs.View,
			Log:  log,
		}
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			var reply StartViewReply
			go func(server int) {
				// fmt.Printf("node-%d sending StartView v=%d to node-%d\n", srv.me, svArgs.View, server)
				srv.peers[server].Call("PBServer_test.StartView", svArgs, &reply)
			}(i)
		}
	}()
}


func (srv *PBServer_test) determineNewViewLog(successReplies []*ViewChangeReply) (
	ok bool, newViewLog []interface{}) {
		// Your code here

	max := 0
	if len(successReplies) >= len(srv.peers) / 2 {
		//decide new one
		ok = true
		//picks the log whose lastest normal view number is the largest.
		for i := 0; i < len(successReplies); i++ {
			if successReplies[i].LastNormalView > max {
				max = successReplies[i].LastNormalView
				//determined the log for the new-view
				newViewLog = successReplies[i].Log
				continue
			}
			if successReplies[i].LastNormalView == max {
				//more than one such logs, it picks the longest log among those
				if (len(successReplies[i].Log) > len(newViewLog)) {
					newViewLog = successReplies[i].Log
				}
			}
		}
	} else {
		ok = false
	}
	return ok, newViewLog
}

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer_test) ViewChange(args *ViewChangeArgs, reply *ViewChangeReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	// checks that the view number included in the message is indeed
	//larger than what it thinks the current view number is.
	if args.View > srv.currentView {
		srv.status = VIEWCHANGE
		reply.LastNormalView = srv.lastNormalView
		reply.Log = srv.log
		reply.Success = true
	} else {
		reply.Success = false
	}

}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer_test) StartView(args *StartViewArgs, reply *StartViewReply) error {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	//the server must again check current view is no bigger than in the RPC message,
	//which would mean that there's been no concurrent view-change for a larger view.
	if srv.currentView > args.View {
		return nil
	} else {
		//sets the new-view as indicated in the message
		//changes its status to be NORMAL
		srv.log = args.Log
		srv.currentView = args.View
		//must set lastnormalview, other wise will return error exit status 1
		srv.lastNormalView = args.View
		srv.status = NORMAL
		srv.commitIndex = len(srv.log) - 1
	}
	return nil
}



func majority(nservers int) int {
	return nservers/2 + 1
}

func Test1ABasicPB(t *testing.T) {
	servers := 3                        //3 servers
	primaryID := GetPrimary(0, servers) //primary ID is determined by view=0
  log.Print("PrimaryID", primaryID)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	for index := 1; index <= 10; index++ {
		xindex := cfg.replicateOne(primaryID, 1000+index, servers) // replicate command 1000+index, expected successful replication to all servers
    log.Print("xindex ?== index",xindex)
    if xindex != index {
			t.Fatalf("got index %v but expected %v", xindex, index)
		}
	}
	fmt.Printf(" ... Passed\n")
}

func Test1AConcurrentPB(t *testing.T) {
	servers := 3                        //3 servers
	primaryID := GetPrimary(0, servers) //primary ID is determined by view=0
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	tries := 5
	for try := 0; try < tries; try++ {
		var wg sync.WaitGroup
		iters := 5
		for i := 0; i < iters; i++ {
			wg.Add(1)
			go func(x int) {
				defer wg.Done()
				val := 2000 + 100*try + x
        vrArgu := common.VrArgu{}
        vrArgu.Op = "test"
        vrArgu.Argu = val
        vrReply := &common.VrReply{}
				// err_test1 := cfg.pbservers[primaryID].Start(vrArgu, vrReply);
        cfg.pbservers[primaryID].Start(vrArgu, vrReply);
        reply,success := vrReply.Reply.(common.TestReply)
        if(!success){
          log.Fatal("convert error in signup")
        }
        if ok := reply.OK; !ok {
					t.Fatalf("node-%d rejected command %v\n", primaryID, val)
				}
			}(i)
		}
		wg.Wait()

		// wait for index (try + 1) * iters to be considered committed
		cfg.waitCommitted(primaryID, (try+1)*iters)

		// check that committed indexes [try*iters, (try+1)*iters] are identical at all servers
		var command interface{}
		for index := 1 + try*iters; index <= (try+1)*iters; index++ {
			cfg.checkCommittedIndex(index, command, majority(servers))
		}
	}
	fmt.Printf(" ... Passed\n")
}

func Test1AFailButCommitPB(t *testing.T) {
	servers := 3 //3 servers
	primaryID := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(primaryID, 3001, servers)

	var wg sync.WaitGroup
	for i := 0; i < 20; i += 4 {
		// disconnect a non-primary server
		cfg.disconnect((primaryID + 1) % servers)

		wg.Add(2)
		go func() {
			defer wg.Done()
			// agree despite replicate disconnected server?
      vrArgu := common.VrArgu{}
      vrArgu.Argu = 3002 + i
      vrReply := &common.VrReply{}
      // err_test1 := cfg.pbservers[primaryID].Start(vrArgu, vrReply);
      cfg.pbservers[primaryID].Start(vrArgu, vrReply);
      reply, _ := vrReply.Reply.(common.TestReply)
			if !reply.OK {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3002+i)
			}
      vrArgu.Argu = 3003 + i
      cfg.pbservers[primaryID].Start(vrArgu, vrReply);
      reply2, _ := vrReply.Reply.(common.TestReply)
			// if _, _, ok := cfg.pbservers[primaryID].Start(3003 + i); !ok {
      if !reply2.OK {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3003+i)
			}
		}()

		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			// re-connect
			cfg.connect((primaryID + 1) % servers)

      vrArgu := common.VrArgu{}
      vrArgu.Argu = 3004 + i
      vrReply := &common.VrReply{}
      //err_test2 := cfg.pbservers[primaryID].Start(vrArgu, vrReply);
      cfg.pbservers[primaryID].Start(vrArgu, vrReply);
      reply_reconn, _ := vrReply.Reply.(common.TestReply)
			if !reply_reconn.OK {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3004+i)
			}
		}()

		wg.Wait()
		cfg.replicateOne(primaryID, 3005, servers)
		// check that all servers replicate the same sequence of commands
		var command interface{}
		for index := 1; index <= 5+i; index++ {
			cfg.checkCommittedIndex(index, command, servers)
		}
		fmt.Printf("iteration i=%d finished\n", i)
	}

	fmt.Printf("  ... Passed\n")
}

func Test1AFailNoCommitPB(t *testing.T) {
	servers := 3 //3 servers
	primaryID := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(primaryID, 4001, servers)

	// disconnect 2 out of 3 servers, both of which are backups
	cfg.disconnect((primaryID + 1) % servers)
	cfg.disconnect((primaryID + 2) % servers)

	// try to replicate command 4002

  vrArgu := common.VrArgu{}
  vrArgu.Op = "test"
  vrArgu.Argu = 4002
  vrReply := &common.VrReply{}
  //err_test_4002 := cfg.pbservers[primaryID].Start(vrArgu, vrReply);
  cfg.pbservers[primaryID].Start(vrArgu, vrReply);
  reply,success := vrReply.Reply.(common.TestReply)
  if(!success){
    log.Fatal("convert error in signup")
  }

  if ok := reply.OK; !ok {
		t.Fatalf("primary rejected the command\n")
	}
  index := reply.Index
	if  index != 2 {
		t.Fatalf("expected index 2, got %v\n", index)
	}
	time.Sleep(2 * time.Second)

	committed := cfg.pbservers[primaryID].IsCommitted(index)
	if committed {
		t.Fatalf("index %d is incorrectly considered to have been committed\n", index)
	}

	// reconnect backups
	cfg.connect((primaryID + 1) % servers)
	cfg.connect((primaryID + 2) % servers)

	cfg.replicateOne(primaryID, 4003, servers)
	index = cfg.replicateOne(primaryID, 4004, servers)

	// disconnect the primary
	cfg.disconnect(primaryID)
  // vrArgu := common.VrArgu{}
  // vrArgu.Op = "test"
  vrArgu.Argu = 4005
  // vrReply := &common.VrReply{}
  //err_test_4005 := cfg.pbservers[primaryID].Start(vrArgu, vrReply);
  cfg.pbservers[primaryID].Start(vrArgu, vrReply);
  reply2,success := vrReply.Reply.(common.TestReply)
  if(!success){
    log.Fatal("convert error in signup")
  }
  index2 := reply2.Index
	// index2, _, ok := cfg.pbservers[primaryID].Start(4005)
	if ok := reply2.OK; !ok {
		t.Fatalf("primary rejected command\n")
	}
	if index2 != (index + 1) {
		t.Fatalf("primary put command at unexpected pos %d\n", index2)
	}
	time.Sleep(2 * time.Second)
	committed = cfg.pbservers[primaryID].IsCommitted(index2)
	if committed {
		t.Fatalf("index %d is incorrectly considered to have been committed\n", index2)
	}

	// reconnect primary
	cfg.connect(primaryID)
	cfg.replicateOne(primaryID, 4006, servers)
	cfg.replicateOne(primaryID, 4007, servers)

	fmt.Printf(" ... Passed\n")
}

func Test1BSimpleViewChange(t *testing.T) {
	servers := 3 //3 servers
	oldPrimary := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(oldPrimary, 5001, servers)
	cfg.checkCommittedIndex(1, 5001, servers)

	// disconnect one backup
	transientBackup := (oldPrimary + 1) % servers
	cfg.disconnect(transientBackup)
	// replicate 5002 at oldPrimary and the remaining connected backup
	cfg.replicateOne(oldPrimary, 5002, majority(servers))
	cfg.checkCommittedIndex(2, 5002, majority(servers))

	// disconnect oldPrimary
	cfg.disconnect(oldPrimary)

	// reconnect the previously disconnected backup
	cfg.connect(transientBackup)

	// change to a new view
	v1 := 1
	cfg.viewChange(v1)
	newPrimary := GetPrimary(v1, servers)

	cfg.replicateOne(newPrimary, 5003, majority(servers))
	cfg.replicateOne(newPrimary, 5004, majority(servers))

	for i := 1; i <= 4; i++ {
		cfg.checkCommittedIndex(i, 5000+i, majority(servers))
	}

	// try to replicate 10 commands 5002...5011 at old disconnected primary

	for i := 0; i < 10; i++ {
    vrArgu := common.VrArgu{}
    vrArgu.Op = "test"
    vrArgu.Argu = 5002 + i
    vrReply := &common.VrReply{}
    //err_test_4002 := cfg.pbservers[oldPrimary].Start(vrArgu, vrReply);
    cfg.pbservers[oldPrimary].Start(vrArgu, vrReply);
    reply,success := vrReply.Reply.(common.TestReply)
    if(!success){
      log.Fatal("convert error in signup")
    }
    ok := reply.OK
		// _, _, ok := cfg.pbservers[oldPrimary].Start(5002 + i)
		if !ok {
			t.Fatalf("old primary %d rejected command\n", oldPrimary)
		}
	}

	// reconnect old primary
	cfg.connect(oldPrimary)

	// replicate 5005 through newPrimary to all 3 servers
	cfg.replicateOne(newPrimary, 5005, servers)
	// check that all 5001...5005 have been replicated at the correct place at all servers
	for i := 1; i <= 5; i++ {
		cfg.checkCommittedIndex(i, 5000+i, servers)
	}
}

func Test1BConcurrentViewChange(t *testing.T) {
	servers := 3 //3 servers
	v0Primary := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(v0Primary, 6001, servers)
	cfg.checkCommittedIndex(1, 6001, servers)

	// disconnect node0
	cfg.disconnect(v0Primary)

	// try to commit command 6002 through disconnected v0Primary, should not succeed
	// test V0 is invalid.
  vrArgu := common.VrArgu{}
  vrArgu.Op = "test"
  vrArgu.Argu = 5999
  vrReply := &common.VrReply{}
  //err_test_5999 := cfg.pbservers[v0Primary].Start(vrArgu, vrReply);
  cfg.pbservers[v0Primary].Start(vrArgu, vrReply);
  reply,success := vrReply.Reply.(common.TestReply)
  if(!success){
    log.Fatal("convert error in signup")
  }
  //
	// index, _, ok := cfg.pbservers[v0Primary].Start(5999)
	if ok := reply.OK; !ok {
		t.Fatalf("primary rejected the command\n")
	}
	if index := reply.Index; index != 2 {
		t.Fatalf("expected index 2, got %v\n", index)
	}
	time.Sleep(2 * time.Second)
	committed := cfg.pbservers[v0Primary].IsCommitted(2)
	if committed {
		t.Fatalf("index 2 is incorrectly considered to have been committed\n")
	}

	// concurrent view change
	// do viewchange
	var wg sync.WaitGroup
	newView := 2
	for v := 1; v <= newView; v++ {
		wg.Add(1)
		go func(view int) {
			defer wg.Done()
			cfg.viewChange(view)
		}(v)
	}
	wg.Wait()

	// reconnect v0Primary
	cfg.connect(v0Primary)

	newView = 5
	for v := 3; v <= newView; v++ {
		wg.Add(1)
		go func(view int) {
			defer wg.Done()
			cfg.viewChange(view)
		}(v)
	}
	wg.Wait()

	newPrimary := GetPrimary(newView, servers)
	cfg.replicateOne(newPrimary, 6002, servers)

	for i := 1; i <= 2; i++ {
		cfg.checkCommittedIndex(i, 6000+i, servers)
	}
}
