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


//
type DB struct{
  like map[string]map[int]bool
  user map[string]User
  msg []common.Msg
  follow map[string]map[string]bool
}

func (db *DB) FollowUser(args *common.FollowUserArgs, reply *common.FollowUserReply) error {
  user := args.User
  follow := args.Follow

  _, ok := db.user[follow]
  if(ok) {
    reply.IsFound = true

    //then followUser
    _, ok := db.follow[user]
    if(ok) {
      //add the follow function
      _, ok := db.follow[user][follow]
      if(ok) {
        //already followU
        reply.IsFollowed = true
        //do unfollow
        delete(db.follow[user], follow)
      }else {
        //unfollow (action:follow)
        reply.IsFollowed = false
        //do follow
        set := make(map[string]bool)
        set[follow] = true
        db.follow[user] = set
      }

    }else {
      //unfollow (action:follow)
      reply.IsFollowed = false
      //do follow
      set := make(map[string]bool)
      set[follow] = true
      db.follow[user] = set
    }


  }else {
    reply.IsFound = false
  }

  return nil
}

func (db *DB) FollowList(args *common.FollowListArgs, reply *common.FollowListReply) error {
  name := args.Name
  _, ok := db.follow[name]
  if(ok) {
    for _, a_msg := range db.msg {
      _, ok = db.follow[name][a_msg.User]
      if(ok) {
        reply.Msg = append(reply.Msg, a_msg)
      }
    }
  }else {
    //do nothing
  }
  return nil
}

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

func (db *DB) DelUser(args *common.DelUserArgs, reply *common.DelUserReply) error {
  name := args.Name
  _, ok := db.user[name]
  if(ok) {
    reply.Success = true
    delete(db.user, name)
  }else {
    reply.Success = false
  }
	return nil
}

func (db *DB) SendMsg(args *common.SendMsgArgs, reply *common.SendMsgReply) error {
  name := args.Name
  value := args.Value
  id := len(db.msg)
  db.msg = append(db.msg, common.Msg{id ,value, name, 0, false})
  reply.Success = true;
	return nil
}

func (db *DB) GetMsg(args *common.GetMsgArgs, reply *common.GetMsgReply) error {
  //name := args.Name
  reply.Msg = db.msg
  reply.Success = true;
	return nil
}

func (db *DB) LikeMsg(args *common.LikeArgs, reply *common.LikeReply) error {
  name := args.Name
  msgid := args.Msgid

  db.msg[msgid].LikeNum += 1
  //add like map if needed
  _, ok := db.like[name]
  if(ok) {
    //append msgid
    db.like[name][msgid] = true
  }else {
    //no like before, add map
    set := make(map[int]bool)
    set[msgid] = true
    db.like[name] = set
  }

  reply.Success = true;
	return nil
}

func (db *DB) UnLikeMsg(args *common.UnLikeArgs, reply *common.UnLikeReply) error {
  name := args.Name
  msgid := args.Msgid

  db.msg[msgid].LikeNum -= 1

  //add like map if needed
  _, ok := db.like[name]
  if(ok) {
    //append msgid
    //test bug
    set := db.like[name]
    delete(set, msgid)
    reply.Success = true
  }else {
    //no like before, it's impossible in unlike
    reply.Success = false
  }
	return nil
}


func (db *DB) LikeList(args *common.LikeListArgs, reply *common.LikeListReply) error {
  name := args.Name
  //add like map if needed
  _, ok := db.like[name]
  if(ok) {
    //append msgid
    //test bug
    reply.Lklist = db.like[name]
    reply.Msg = db.msg
    reply.Success = true
  }else {
    //no like before, it's impossible in unlike
    reply.Success = false
  }
	return nil
}

func (db *DB) IsLike(args *common.IsLikeArgs, reply *common.IsLikeReply) error {
  name := args.Name
  msgid := args.Msgid
  _, ok := db.like[name][msgid]
  if(ok) {
    reply.Success = true
  }else {
    reply.Success = false
  }
  return nil
}


type vrCode struct {
  db *DB
  ipaddress string
}

func (vc *vrCode) run(args *common.VrArgs, reply *common.VrReply) error {
  switch(args.op){

  }
  vc.db.
}


func main(){
  db := new(DB)
  
  vc = vrCode(db,"000,000,,")

  rpc.Register(vc)
  rpc.HandleHTTP()
  l, e := net.Listen("tcp", ":8081")
  if e != nil {
  	log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)
}
