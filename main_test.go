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

func TestSignup(){
  
}
