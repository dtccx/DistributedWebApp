package main

import (
  "testing"
  "log"
)


func TestLogin(t *testing.T){
  // var user map[string]User
  user = make(map[string]User)
  var temp = User{"name", "password"}
  user["name"] = temp
  log.Println(user)
  ret := _login("name", "password")
  if ret != "true"{
    t.Fatalf("TestLogin fail")
  }
}

func TestisLike(t *testing.T) {
  // msg = {
  // 			  {
  //           ID  : 0,
  //           Value : "I like debuging :)",
  //           User  : "usera",
  //           LikeNum  : 2,
  //           IsLiked : false,
  // 			  },
  //         {
  //           ID  : 1,
  //           Value : "I literally like debuging :)",
  //           User  : "userb",
  //           LikeNum  : 3,
  //           IsLiked : false,
  // 			  },
  //         {
  //           ID  : 2,
  //           Value : "I really like debuging :)",
  //           User  : "userc",
  //           LikeNum  : 3,
  //           IsLiked : false,
  // 			  },
  //  			}
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
}
