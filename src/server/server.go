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
  "sync"
  "vrproxy"
  "encoding/gob"
)

type Arith struct {
  client *rpc.Client
}
var arith *Arith
var vp *vrproxy.VrProxy


var store = sessions.NewCookieStore([]byte("something-very-secret"))


func main() {
  // gob.RegisterName(common.SignArgs{})
  // gob.RegisterName("SignArgs",common.SignArgs{})
  // gob.Register(common.SignReply{})
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

  args := &common.FollowListArgs{name}
  var reply common.FollowListReply
  err := arith.client.Call("DB.FollowList", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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
  args := &common.FollowUserArgs{user, follow}
  var reply common.FollowUserReply
  err := arith.client.Call("DB.FollowUser", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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
    var mu sync.Mutex
    mu.Lock()
    var err error
    args := &common.LogArgs{name}
    var reply common.LogReply
    log.Println("client", t.client)
    go func() {
      err = t.client.Call("DB.Login", args, &reply)
      mu.Unlock()
      //client.Call("DB.Login", args, &reply)
    }()
    mu.Lock()
    if err != nil {
      log.Fatal("arith error:", err)
    }
    if(reply.Success && reply.Password == password){
      log.Println(reply.Password)
      mu.Unlock()
      return "true"
      }else {
        mu.Unlock()
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

    go func() {
      // var args interface{} = common.SignArgs{name, password}
      args := common.SignArgs{name, password}
      // gob.RegisterName("common.SignArgs", args)
      gob.Register(common.SignArgs{})
      // log.Println(args)
      // _args, ok := args.(interface{})
      // if(!ok){
      //   fmt.Println("fail to convert args to interface{}")
      // }
      vrArgu := &common.VrArgu{}
      vrArgu.Argu = args
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
      }
    }()

  } else {
    t, _ := template.ParseFiles("index.html")
    log.Println(t.Execute(w, nil))
  }
}

func delUser(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  var temp interface{} = "user"
  name := session.Values[temp].(string)
  args := &common.DelUserArgs{name}
  var reply common.DelUserReply
  err := arith.client.Call("DB.DelUser", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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
    args := &common.SendMsgArgs{name, value}
    var reply common.SendMsgReply
    err := arith.client.Call("DB.SendMsg", args, &reply)

    if err != nil {
      log.Fatal("arith error:", err)
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


  args := &common.GetMsgArgs{name}
  var reply common.GetMsgReply
  err := arith.client.Call("DB.GetMsg", args, &reply)

  if err != nil {
    log.Fatal("arith error:", err)
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
  args := &common.IsLikeArgs{user, msgid}
  var reply common.IsLikeReply
  err := arith.client.Call("DB.IsLike", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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

  args := &common.LikeArgs{name, msgid}
  var reply common.LikeReply
  err := arith.client.Call("DB.LikeMsg", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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

  args := &common.UnLikeArgs{name, msgid}
  var reply common.UnLikeReply
  err := arith.client.Call("DB.UnLikeMsg", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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

  args := &common.LikeListArgs{name}
  var reply common.LikeListReply
  err := arith.client.Call("DB.LikeList", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
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
