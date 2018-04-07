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
  log.Fatal(http.ListenAndServe(":8080", nil))
}




func MyHandler(w http.ResponseWriter, r *http.Request){
  // Get a session. We're ignoring the error resulted from decoding an
  // existing session: Get() always returns a session, even if empty.
  session, _ := store.Get(r, "session-name")
  // Set some session values.
  session.Values["foo"] = "bar"
  session.Values[42] = 43
  // Save it before we write to the response/return from the handler.
  session.Save(r, w)
}

func login(w http.ResponseWriter, r *http.Request) {
    //name := r.FormValue("name")
    name := r.FormValue("user")
    password := r.FormValue("password")
    i, ok := user[name]
    if(ok && i.Password == password){
      log.Println("log success")
      session, _ := store.Get(r, "session-name")
      // Set some session values.
      session.Values["authenticated"] = true
      session.Values["user"] = name
      // Save it before we write to the response/return from the handler.
      session.Save(r, w)
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

func sendMsg(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    t, _ := template.ParseFiles("index.html")
    log.Println(t.Execute(w, nil))
  } else {
    value := r.FormValue("value")
    id := len(msg) - 1
    //session, _ := store.Get(r, "session-name")
    //name := session.Values["user"]
    name := "sb"
    msg = append(msg, Msg{id ,value, name, 0})
    log.Println(msg)
  }
}

// func getMsg() {
//   if r.Method == "GET" {
//     return msg
//   } else {
//     t, _ := template.ParseFiles("index.html")
//     log.Println(t.Execute(w, nil))
//   }
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
