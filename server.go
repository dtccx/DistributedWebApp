package  main

import (
  "fmt"
  "net/http"
  "log"
  "html/template"
  "net/rpc"
  "./src/common"
  "github.com/gorilla/sessions"
  "encoding/json"
  "strconv"
)

type Arith struct {
  client *rpc.Client
}
var arith *Arith

var like map[string]map[int]bool
var user map[string]User
var msg []Msg
//list msg
//var msg map[int]Msg

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type User struct {
	Name          string
	Password      string
}

type Msg struct {
  ID            int
  Value         string
  User          string
  LikeNum       int
  IsLiked       bool
}


func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  client, err := rpc.DialHTTP("tcp", "localhost:8081")
  if err != nil {
    log.Fatal("dialing:", err)
  }

  arith = &Arith{client: client}

  http.HandleFunc("/User/Login", login)
  http.HandleFunc("/User/Register", signup)
  http.HandleFunc("/SendMsg", sendMsg)
  http.HandleFunc("/GetMsg", getMsg)
  http.HandleFunc("/DelUser", delUser)
  http.HandleFunc("/LikeMsg", likeMsg)
  http.HandleFunc("/UnlikeMsg", unlikeMsg)
  http.HandleFunc("/LikeList", likeList)
  http.ListenAndServe(":8080", nil)
  // Tries to connect to localhost:1234 using HTTP protocol (The port on which rpc server is listening)


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

  args := &common.LogArgs{name}
  var reply common.LogReply
  log.Println("client", t.client)
  err := t.client.Call("DB.Login", args, &reply)
  //client.Call("DB.Login", args, &reply)
  if err != nil {
    log.Fatal("arith error:", err)
  }
  if(reply.Success && reply.Password == password){
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
    log.Print(name)
    log.Print(password)
    args := &common.SignArgs{name, password}
    var reply common.SignReply
    err := arith.client.Call("DB.Signup", args, &reply)
    //client.Call("DB.Login", args, &reply)
    if err != nil {
      log.Fatal("arith error:", err)
    }

    if(!reply.Success){
      //user exsit
      log.Println("User already exist")
      fmt.Fprintf(w, "0") //exsit

    }else {
      //user[name] = User{name, password}
      log.Print("map:", user)
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


  log.Println(user)
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
      log.Println("User already exist")
      fmt.Fprintf(w, "0") //exsit

    }else {
      //user[name] = User{name, password}
      log.Print("map:", user)
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
  _, ok := like[user][msgid]
  if(ok) {
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

  args := common.LikeArgs{name, msgid}
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
  msg[msgid].LikeNum -= 1

  //add like map if needed
  _, ok := like[name]
  if(ok) {
    //append msgid
    set := like[name]
    delete(set, msgid)
  }else {
    //no like before, it's impossible in unlike

  }
}

func likeList(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  log.Println(session)
  var temp interface{} = "user"
  name := session.Values[temp].(string)

  //send the list of liking msgid to client
  var like_list []int
  set, ok := like[name]
  if(ok) {
    for k, _ := range set {
      like_list = append(like_list, k)
    }
  }else {
    //no possilbe

  }
  //json
  var msg_get []Msg
  for _, i := range like_list {
    msg_get = append(msg_get, msg[i])
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
