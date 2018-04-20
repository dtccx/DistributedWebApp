package common

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
