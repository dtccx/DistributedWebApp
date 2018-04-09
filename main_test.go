package main

import (
  "testing"
  // "log"
  "fmt"
  "github.com/gorilla/sessions"
  "net/http/httptest"
  "io/ioutil"
  // "strings"
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


func TestLogin2(t *testing.T){
  store = sessions.NewCookieStore([]byte("something-very-secret"))
  user = make(map[string]User)
  user["user"] = User{"user", "password"}
  data := url.Values{}
  data.Set("user", "user")
  data.Add("password", "password")

  r := httptest.NewRequest("GET", "/User/Login?"+data.Encode(), nil)
  w := httptest.NewRecorder()
  login(w,r)

  resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

  ret := string(body)
  if ret=="false"{
    t.Fatalf("TestLogin2 fail")
  }
}

func TestisLike(t *testing.T) {
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
    }
  }
  fmt.Printf("Liked Passed\n")
}
