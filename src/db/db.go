package main

import (
  "errors"
  "log"
  "net/rpc"
  "net"
  "net/http"
  "common"
)

type Arith int

func (t *Arith) Multiply(args *common.Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *common.Args, quo *common.Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main(){
  arith := new(Arith)
  rpc.Register(arith)
  rpc.HandleHTTP()
  l, e := net.Listen("tcp", ":1234")
  if e != nil {
  	log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)
}
