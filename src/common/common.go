package common

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

type msgArgs struct {
  Name string
}
