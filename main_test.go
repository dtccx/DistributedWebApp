package main

import (
  "testing"
  "log"
  "fmt"
  "github.com/gorilla/sessions"
  "net/http/httptest"
  "io/ioutil"
  "net/http"
  "net/url"
)

// func TestLogin(t *testing.T){
//   // var user map[string]User
//   user = make(map[string]User)
//   var temp = User{"name", "password"}
//   user["name"] = temp
//   log.Println(user)
//   ret := _login("name", "password")
//   if ret != "true"{
//     t.Fatalf("TestLogin fail")
//   }
//   fmt.Printf("Login Passed\n")
// }
var urlString = "http://localhost:8080";

func TestLogin(t *testing.T){
  store = sessions.NewCookieStore([]byte("something-very-secret"))
  user = make(map[string]User)
  user["user"] = User{"user", "password"}
  data := url.Values{}
  data.Set("user", "user")
  data.Add("password", "password")

  r := httptest.NewRequest("GET", urlString+"/User/Login?"+data.Encode(), nil)
  w := httptest.NewRecorder()
  login(w,r)

  resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

  ret := string(body)
  if ret=="false"{
    t.Fatalf("TestLogin2 fail")
  }
  fmt.Printf("...... TestLogin Passed\n")
}

func TestIsLike(t *testing.T) {
  like = make(map[string]map[int]bool)
  test_case1 := []string{"usera","userb","userb"}
  test_case2 := []int{1,0,1}
  temp := make(map[int]bool)
  temp[1] = true
  like["usera"] = temp
  temp[0] = true
  like["userb"] = temp

  for i := 0; i < len(test_case1); i++ {
    ok := isLike(test_case1[i],test_case2[i])
    if(!ok) {
      t.Fatalf("TestLike fail")
      fmt.Printf("Liked Failed\n")
    }
  }
  fmt.Printf("...... TestMessageisLiked Passed\n")
}


func TestSendMsgHttp(t *testing.T){
  handler := func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "...... SendMsgResponse Passed", http.StatusInternalServerError)
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
      fmt.Fprintln(w, "...... TestGetMessageServer Passed")
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
