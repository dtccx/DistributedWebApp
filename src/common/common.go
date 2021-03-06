package common

// var Address = []string{":8081",":8082",":8083"}


  // Address[0] = ":8081"
  // Address[1] = ":8082"
  // Address[2] = ":8083"

type DealPrimayArgs struct {

}

type DealPrimayReply struct {
  OK            bool
}

type Msg struct {
  ID            int
  Value         string
  User          string
  LikeNum       int
  IsLiked       bool
}

type LogArgs struct {
	Name string
}

type LogReply struct {
  Password string
	Success bool
}

type SignArgs struct {
  Name string
  Password string
}

type SignReply struct {
  Success bool
}

type DelUserArgs struct {
	Name string
}

type DelUserReply struct {
	Success bool
}

type SendMsgArgs struct {
  Name string
	Value string
}

type SendMsgReply struct {
  Success bool
}

type GetMsgArgs struct {
  Name string
}

type GetMsgReply struct {
	Msg 		[]Msg
  Success bool
}

type LikeArgs struct {
	Name		string
	Msgid		int
}
type LikeReply struct {
	Success bool
}

type UnLikeArgs struct {
	Name		string
	Msgid		int
}
type UnLikeReply struct {
	Success bool
}

type LikeListArgs struct {
	Name		string
}
type LikeListReply struct {
	Lklist	map[int]bool
	Msg			[]Msg
	Success bool
}

type IsLikeArgs struct {
	Name 		string
	Msgid		int
}

type IsLikeReply struct {
	Success bool
}

type FollowUserArgs struct {
  User    string
  Follow  string
}
type FollowUserReply struct {
  IsFound     bool
  IsFollowed  bool
}

type FollowListArgs struct {
  Name    string
}

type FollowListReply struct {
  Msg     []Msg
}

type VrArgu struct{
  Op  string
  Argu  interface{}
}

type VrReply struct{
  Reply interface{}
}

type VrViewChangeArgu struct{
  View  int
}

type VrViewChangeReply struct{

}
