package main

import (
  "fmt"
  "net/http"
  "log"
  "html/template"
  "github.com/gorilla/sessions"
)

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
}


func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  user = make(map[string]User)
  // http.HandleFunc("/", HomePage)
  http.HandleFunc("/User/Login", login)
  http.HandleFunc("/User/Register", signup)
  http.HandleFunc("/SendMsg", sendMsg)
  // http.HandleFunc("/GetMsg", getMsg)
  log.Fatal(http.ListenAndServe(":8080", nil))
}



func login(w http.ResponseWriter, r *http.Request) {
    //name := r.FormValue("name")
    name := r.FormValue("user")
    password := r.FormValue("password")
    i, ok := user[name]
    if(ok && i.Password == password){
      log.Println("log success")
      session, _ := store.Get(r, "user_session")
      // Set some session values.
      //session.Values["authenticated"] = true
      var temp interface{} = "user"
      session.Values[temp] = name
      // Save it before we write to the response/return from the handler.
      session.Save(r, w)
      log.Print(session)
      fmt.Fprintf(w, "true")
    }else {
      //log fails
      log.Println("log fails")
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
    _, ok := user[name]
    if(ok){
      //user exsit
      log.Println("User already exist")
      fmt.Fprintf(w, "0") //exsit

    }else {
      user[name] = User{name, password}
      log.Print("map:", user)
    }

  } else {
    t, _ := template.ParseFiles("index.html")
    log.Println(t.Execute(w, nil))
  }
}

func msgdel(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "user_session")
  var temp interface{} = "user"
  name := session.Values[temp].(string)
  _, ok := user[name]
  if(ok) {
    delete(user, name)
  }
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
    value := r.FormValue("value")
    id := len(msg)
    session, _ := store.Get(r, "user_session")
    log.Println(session)
    var temp interface{} = "user"
    name := session.Values[temp].(string)
    //name := "sb"
    msg = append(msg, Msg{id ,value, name, 0})
    log.Println(msg)
    fmt.Fprintf(w, value)

}

// func getMsg() {
//   for(i := len(msg) - 1; i >= 0; i--){
//     msg[i].Value
//   }
// }

// func like(w http.ResponseWriter, r *http.Request){
//   name := //user Name
//   msgid :=
//   msg[msgid].LikeNum += 1
//
//   //add like map if needed
//
// }

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
