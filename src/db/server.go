package  main

import (
  "fmt"
  "net/http"
  "log"
  "html/template"
  "net/rpc"
  "common"
  "github.com/gorilla/sessions"
  "encoding/json"
  "strconv"
  "vrproxy"
  "encoding/gob"
  "net"
)

type Arith struct {
  client *rpc.Client
}
var arith *Arith
var vp *vrproxy.VrProxy


var store = sessions.NewCookieStore([]byte("something-very-secret"))

func createServer(clients []*rpc.Client, serverIndex int) (*rpc.Client, *PBServer){
  ps := Make(clients, serverIndex, 0)
  server := rpc.NewServer()
  server.Register(ps)
  port := ":"+strconv.Itoa(8081+serverIndex)
  l,listenError := net.Listen("tcp", port)
  if(listenError!=nil){
    log.Println(listenError)
  }
  go server.Accept(l)
  client, err := rpc.Dial("tcp", port)
  clients[0] = client
  if(err!=nil){
    log.Println(err)
  }
  log.Println(client==nil)

  clients[serverIndex] = client

  return client, ps
}

func main() {
  serverNum := 3
  clients := make([]*rpc.Client, serverNum)
  for i := 0; i < serverNum; i++ {
	   createServer(clients, i)
	}

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

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  client, err := rpc.Dial("tcp", "localhost:8081")
  if err != nil {
    log.Fatal("dialing:", err)
  }
  vp = vrproxy.Make(client)

  arith = &Arith{client: client}

  http.HandleFunc("/User/Login", login)
  http.HandleFunc("/User/Register", signup)
  http.HandleFunc("/SendMsg", sendMsg)
  http.HandleFunc("/GetMsg", getMsg)
  http.HandleFunc("/DelUser", delUser)
  http.HandleFunc("/LikeMsg", likeMsg)
  http.HandleFunc("/UnlikeMsg", unlikeMsg)
  http.HandleFunc("/LikeList", likeList)
  http.HandleFunc("/FollowUser", followUser)
  http.HandleFunc("/FollowList", followList)
  http.ListenAndServe(":8080", nil)
  // log.Print("here")
  // Tries to connect to localhost:1234 using HTTP protocol (The port on which rpc server is listening)
}

func followList(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  var temp interface{} = "user"
  name := session.Values[temp].(string)

  vrArgu := &common.VrArgu{}
  args := common.FollowListArgs{name}
  vrArgu.Argu = args
  vrArgu.Op = "DB.FollowList"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.FollowListReply)
  if(!ok){
    log.Fatal("convert error in FollowList")
  }

  j, _ := json.Marshal(reply.Msg)
  fmt.Fprintf(w, string(j))
  log.Println(string(j))

}


func followUser(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  var temp interface{} = "user"
  user := session.Values[temp].(string)
  follow := r.FormValue("user")

  vrArgu := &common.VrArgu{}
  args := common.FollowUserArgs{user, follow}
  vrArgu.Argu = args
  vrArgu.Op = "DB.FollowUser"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.FollowUserReply)
  if(!ok){
    log.Fatal("convert error in FollowUser")
  }


  if(reply.IsFound){
    //user exsit
    log.Println("User valid")

    if(reply.IsFollowed) {
      //do unfollow
      fmt.Fprintf(w, "10")
    }else {
      //follow
      fmt.Fprintf(w, "11")
    }

  }else {
    //user[name] = User{name, password}
    log.Print("Not valid user")
    fmt.Fprintf(w, "00")
  }
}



func login(w http.ResponseWriter, r *http.Request) {
    //name := r.FormValue("name")
    name := r.FormValue("user")
    password := r.FormValue("password")

    // log.Println("login name:", name)
    // log.Println("login password:", password)

    ret := arith._login(name, password)
    if ret == "true" {
      session, _ := store.Get(r, "user_session")
      // Set some session values.
      //session.Values["authenticated"] = true
      var temp interface{} = "user"
      session.Values[temp] = name
      // Save it before we write to the response/return from the handler.
      session.Save(r, w)
      log.Print(session)
    }
    fmt.Fprintf(w, ret)
}

func(t *Arith) _login(name string, password string) string{
    // log.Println("login name:", name)
    // log.Println("login password:", password)
    args := &common.LogArgs{name}
    //var reply common.LogReply

    vrArgu := &common.VrArgu{}
    vrArgu.Argu = args
    vrArgu.Op = "DB.Login"
    vrReply := &common.VrReply{}
    vrErr := vp.CallVr(vrArgu, vrReply)
    if(vrErr!=nil){
      log.Fatal("vr error:", vrErr)
    }
    reply,ok := vrReply.Reply.(common.LogReply)
    if(!ok){
      log.Fatal("convert error in login")
    }

    // log.Println("client", t.client)
    if(reply.Success && reply.Password == password){
      log.Println(reply.Password)
      return "true"
      }else {
        return "false"
      }

}

func signup(w http.ResponseWriter, r *http.Request) {
  fmt.Println("method:", r.Method)
  if r.Method == "GET" {
    r.ParseForm()
    fmt.Println(r.Form)
    name := r.FormValue("user")
    password := r.FormValue("password")

      args := common.SignArgs{name, password}
      gob.Register(common.SignArgs{})
      vrArgu := &common.VrArgu{}
      vrArgu.Argu = args
      vrArgu.Op = "DB.Signup"
      vrReply := &common.VrReply{}
      vrErr := vp.CallVr(vrArgu, vrReply)
      if(vrErr!=nil){
        log.Fatal("vr error:", vrErr)
      }
      reply,ok := vrReply.Reply.(common.SignReply)
      if(!ok){
        log.Fatal("convert error in signup")
      }
      if(!reply.Success){
        //user exsit
        log.Println("User already exist")
        fmt.Fprintf(w, "0") //exsit

      }else {
        //user[name] = User{name, password}
        log.Print("signup success")
        fmt.Fprintf(w, "1") //exsit
      }

  } else {
    t, _ := template.ParseFiles("index.html")
    log.Println(t.Execute(w, nil))
  }
}

func delUser(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  var temp interface{} = "user"
  name := session.Values[temp].(string)

  vrArgu := &common.VrArgu{}
  args := common.DelUserArgs{name}
  vrArgu.Argu = args
  vrArgu.Op = "DB.DelUser"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
    // if need to refresh
    // t, _ := template.ParseFiles("index.html")
    // log.Println(t.Execute(w, nil))
  }
  reply,ok := vrReply.Reply.(common.DelUserReply)
  if(!ok){
    log.Fatal("convert error in login")
  }

  if(!reply.Success){
    //user not exist
    log.Println("User not exist")
    fmt.Fprintf(w, "0") //exsit
  }else {
    log.Print("delete successfully")
  }

  fmt.Fprintf(w, "0")
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
    value := r.FormValue("value")
    //id := len(msg)
    session, _ := store.Get(r, "user_session")
    log.Println(session)
    var temp interface{} = "user"
    name := session.Values[temp].(string)

    vrArgu := &common.VrArgu{}
    args := common.SendMsgArgs{name, value}
    vrArgu.Argu = args
    vrArgu.Op = "DB.SendMsg"
    vrReply := &common.VrReply{}
    vrErr := vp.CallVr(vrArgu, vrReply)
    if(vrErr!=nil){
      log.Fatal("vr error:", vrErr)
    }
    reply,ok := vrReply.Reply.(common.SendMsgReply)
    if(!ok){
      log.Fatal("convert error in login")
    }


    if(!reply.Success){
      //user exsit
      log.Println("send fail")
      fmt.Fprintf(w, "0") //exsit

    }else {
      //user[name] = User{name, password}
      log.Print("send success")
    }

    fmt.Fprintf(w, value)
}



func getMsg(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  log.Println(session)
  var temp interface{} = "user"
  name := session.Values[temp].(string)
  index_str := r.FormValue("index")
  index, _ := strconv.Atoi(index_str) //_, error
  var msg_get []common.Msg
  var msg []common.Msg
  msgnum := 3

  vrArgu := &common.VrArgu{}
  args := common.GetMsgArgs{name}
  vrArgu.Argu = args
  vrArgu.Op = "DB.GetMsg"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.GetMsgReply)
  if(!ok){
    log.Fatal("convert error in getMsg")
  }

  if(!reply.Success){
    //user exsit
    log.Println("User already exist")
    return

  }else {
    //user[name] = User{name, password}
    msg = reply.Msg
  }



  //everytime the num of sending
  if(index == -2){
    for i := len(msg) - 1; i >= len(msg) - msgnum && i >= 0; i-- {
      //i is id of the msg
      temp := msg[i]
      if(isLike(name, i)){
        temp.IsLiked = true
      }else {
        temp.IsLiked = false
      }
      msg_get = append(msg_get, temp)
    }
  }else {
    for i := index; i >= index - msgnum + 1 && i >= 0; i-- {
      temp := msg[i]
      if(isLike(name, i)){
        temp.IsLiked = true
      }else {
        temp.IsLiked = false
      }
      msg_get = append(msg_get, temp)
    }
  }
  log.Println(msg_get)
  j, _ := json.Marshal(msg_get)
  fmt.Fprintf(w, string(j))
  log.Println(string(j))
  // var f interface{}
  // log.Println(json.Unmarshal([]byte(string(j)), &f))
  // log.Println(f)
}

func isLike(user string, msgid int) bool {
  vrArgu := &common.VrArgu{}
  args := common.IsLikeArgs{user, msgid}
  vrArgu.Argu = args
  vrArgu.Op = "DB.IsLike"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.IsLikeReply)
  if(!ok){
    log.Fatal("convert error in isLike")
  }

  if(reply.Success){
    return true
  }else {
    return false
  }

}

func likeMsg(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "user_session")
  log.Println(session)
  var temp interface{} = "user"
  name := session.Values[temp].(string)
  msgid_str := r.FormValue("msgid")
  msgid, _ := strconv.Atoi(msgid_str)

  vrArgu := &common.VrArgu{}
  args := common.LikeArgs{name, msgid}
  vrArgu.Argu = args
  vrArgu.Op = "DB.LikeMsg"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.LikeReply)
  if(!ok){
    log.Fatal("convert error in LikeMsg")
  }

  if(!reply.Success){
    //user exsit
    log.Println("like fail")
    return

  }else {
    //user[name] = User{name, password}
    log.Println("like success")
    return
  }


}

func unlikeMsg(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  log.Println(session)
  var temp interface{} = "user"
  name := session.Values[temp].(string)

  msgid_str := r.FormValue("msgid")
  msgid, _ := strconv.Atoi(msgid_str)

  vrArgu := &common.VrArgu{}
  args := common.UnLikeArgs{name, msgid}
  vrArgu.Argu = args
  vrArgu.Op = "DB.UnLikeMsg"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.UnLikeReply)
  if(!ok){
    log.Fatal("convert error in UnLikeMsg")
  }

  if(!reply.Success){
    //user exsit
    log.Println("unlike fail")
    return

  }else {
    //user[name] = User{name, password}
    log.Println("unlike success")
    return
  }

}

func likeList(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  log.Println(session)
  var temp interface{} = "user"
  name := session.Values[temp].(string)

  vrArgu := &common.VrArgu{}
  args := common.LikeListArgs{name}
  vrArgu.Argu = args
  vrArgu.Op = "DB.LikeList"
  vrReply := &common.VrReply{}
  vrErr := vp.CallVr(vrArgu, vrReply)
  if(vrErr!=nil){
    log.Fatal("vr error:", vrErr)
  }
  reply,ok := vrReply.Reply.(common.LikeListReply)
  if(!ok){
    log.Fatal("convert error in LikeList")
  }

  if(!reply.Success){
    //user exsit
    log.Println("likelist fail")
    return

  }

  //send the list of liking msgid to client
  var like_list []int
  set := reply.Lklist

  for k, _ := range set {
    like_list = append(like_list, k)
  }

  //json
  var msg_get []common.Msg
  for _, i := range like_list {
    msg_get = append(msg_get, reply.Msg[i])
  }
  j, _ := json.Marshal(msg_get)
  fmt.Fprintf(w, string(j))
  log.Println(string(j))

}

func HomePage(w http.ResponseWriter, r *http.Request){
    t, err := template.ParseFiles("index.html") //parse the html file homepage.html
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, nil) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}
