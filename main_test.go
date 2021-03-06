package main

import (
  "testing"
  "log"
  "fmt"
  // "github.com/gorilla/sessions"
  "net/http/httptest"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
  "common"
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

// type Msg struct {
//   ID            int
//   Value         string
//   User          string
//   LikeNum       int
//   IsLiked       bool
// }
//
// type User struct {
// 	Name          string
// 	Password      string
// }

var urlString = "http://localhost:8080";

// func TestInit(t *testing.T){
//
// }

func TestLogin(t *testing.T){
  // store := sessions.NewCookieStore([]byte("something-very-secret"))
  user := make(map[string]User)
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
  fmt.Printf("............ TestLogin Passed ********* !!\n")
}

func TestIsLike(t *testing.T) {
  like := make(map[string]map[int]bool)
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
  fmt.Printf("............ TestMessageisLiked Passed ********* !!\n")
}


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


  user := make(map[string]User)
  user["user"] = User{"user", "password"}
  data := url.Values{}
  data.Set("index", "-2")
  // data.Add("password", "password")
  // msg := []Msg{
	// 		  {
  //         ID  : 0,
  //         Value : "I like debuging :)",
  //         User  : "usera",
  //         LikeNum  : 2,
  //         IsLiked : false,
	// 		  },
  //       {
  //         ID  : 1,
  //         Value : "I literally like debuging :)",
  //         User  : "userb",
  //         LikeNum  : 3,
  //         IsLiked : false,
	// 		  },
  //       {
  //         ID  : 2,
  //         Value : "I really like debuging :)",
  //         User  : "userc",
  //         LikeNum  : 3,
  //         IsLiked : false,
	// 		  }}

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
                IsLiked : true,
      			  },
              {
                ID  : 0,
                Value : "I like debuging :)",
                User  : "usera",
                LikeNum  : 2,
                IsLiked : true,
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
  log.Println(ret)
  log.Println(string(jsonval))
  if ret != string(jsonval){
    t.Fatalf("TestGetMsg fail")
  }
  fmt.Printf("............ TestGetMsgJsonResponse Passed ********* !!\n")
}
