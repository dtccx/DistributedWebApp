import (
  "testing"
  "main"
)


func TestLogin(t *testing.T){
  // var user map[string]User
  user["name"] = User{"name", "password"}
  ret := _login("name", "password")
  if ret!="true"{
    t.Fatalf("TestLogin fail")
  }
}
