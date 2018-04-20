package main

import (
  // "errors"
  "log"
  "net/rpc"
  "net"
  "net/http"
  "common"
)

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
//
type DB struct{
  like map[string]map[int]bool
  user map[string]User
  msg []Msg
}

// func (db *DB) Multiply(args *common.logArgs, reply *int) error {
// 	*reply = args.A * args.B
// 	return nil
// }
//
func (db *DB) Login(args *common.LogArgs, reply *common.LogReply) error {
  name := args.Name
  i, ok := db.user[name]
  if(ok){
    reply.Password = i.Password
    reply.Success = true
  }else{
    reply.Success = false
  }
	return nil
}

func (db *DB) Signup(args *common.SignArgs, reply *common.SignReply) error {
  log.Println(args.Password)
  name := args.Name
  password := args.Password
  _, ok := db.user[name]
  if(ok){
    reply.Success = false
  }else{
    reply.Success = true
    db.user[name] = User{name, password}
  }
	return nil
}

func (db *DB) SendMsg(args *common.SendMsgArgs, reply *common.SendMsgReply) error {
  name := args.Name
  value := args.Value
  id := len(db.msg)
  db.msg = append(db.msg, Msg{id ,value, name, 0, false})
  reply.Success = true;
	return nil
}


// func (t *DB) Divide(args *common.Args, quo *common.Quotient) error {
// 	if args.B == 0 {
// 		return errors.New("divide by zero")
// 	}
// 	quo.Quo = args.A / args.B
// 	quo.Rem = args.A % args.B
// 	return nil
// }


func main(){
  db := new(DB)
  db.user = make(map[string]User)
  db.like = make(map[string]map[int]bool)
  rpc.Register(db)
  rpc.HandleHTTP()
  l, e := net.Listen("tcp", ":8081")
  if e != nil {
  	log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)
}
