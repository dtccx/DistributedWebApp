package main

import (
	"sync"
	"log"
  "net/rpc"
	// "net"
	// "flag"
	"common"
	"encoding/gob"
	// "os"
	// "os/signal"
	"errors"
	// "net/http"
)

// the 3 possible server status
const (
	NORMAL = iota  	//0
	VIEWCHANGE			//1
	RECOVERING			//2
)

// PBServer defines the state of a replica server (either primary or backup)
type PBServer struct {
	mu             sync.Mutex          // Lock to protect shared access to this peer's state
	peers          []*rpc.Client // RPC end points of all peers
	me             int                 // this peer's index into peers[]
	currentView    int                 // what this peer believes to be the current active view
	status         int                 // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int                 // the latest view which had a NORMAL status

	log         []interface{} // the log of "commands"
	commitIndex int           // all log entries <= commitIndex are considered to have been committed.

	db					*DB
}

// Prepare defines the arguments for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC args struct
type PrepareArgs struct {
	View          int         // the primary's current view
	PrimaryCommit int         // the primary's commitIndex
	Index         int         // the index position at which the log entry is to be replicated on backups
	Entry         interface{} // the log entry to be replicated
}

// PrepareReply defines the reply for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC reply struct
type PrepareReply struct {
	View    int  // the backup's current view
	Success bool // whether the Prepare request has been accepted or rejected
}

// RecoverArgs defined the arguments for the Recovery RPC
type RecoveryArgs struct {
	View   int // the view that the backup would like to synchronize with
	Server int // the server sending the Recovery RPC (for debugging)
}

type RecoveryReply struct {
	View          int           // the view of the primary
	Entries       []interface{} // the primary's log including entries replicated up to and including the view.
	PrimaryCommit int           // the primary's commitIndex
	Success       bool          // whether the Recovery request has been accepted or rejected
}

type ViewChangeArgs struct {
	View int // the new view to be changed into
}

type ViewChangeReply struct {
	LastNormalView int           // the latest view which had a NORMAL status at the server
	Log            []interface{} // the log at the server
	Success        bool          // whether the ViewChange request has been accepted/rejected
}

type StartViewArgs struct {
	View int           // the new view which has completed view-change
	Log  []interface{} // the log associated with the new new
}

type StartViewReply struct {
}


func GetPrimary(view int, nservers int) int {
	return view % nservers
}

func (srv *PBServer) IsCommitted(index int) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.commitIndex >= index {
		return true
	}
	return false
}

func (srv *PBServer) ViewStatus() (currentView int, statusIsNormal bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView, srv.status == NORMAL
}

func (srv *PBServer) GetEntryAtIndex(index int) (ok bool, command interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.log) > index {
		return true, srv.log[index]
	}
	return false, command
}

func (srv *PBServer) Kill() {
	// Your code here, if necessary
}

func Make(peers []*rpc.Client, me int, startingView int) *PBServer {
	gob.Register(common.SignArgs{})
	gob.Register(common.SignReply{})
	gob.Register(common.LogArgs{})
	gob.Register(common.DelUserArgs{})

	srv := &PBServer{
		peers:          peers,
		me:             me,
		currentView:    startingView,
		lastNormalView: startingView,
		status:         NORMAL,
		db:							new(DB),
	}
	srv.db.user = make(map[string]User)
	srv.db.like = make(map[string]map[int]bool)
	srv.db.follow = make(map[string]map[string]bool)
	// all servers' log are initialized with a dummy command at index 0
	var v interface{}
	srv.log = append(srv.log, v)

	// Your other initialization code here, if there's any
	return srv
}

func (srv *PBServer) Start(args common.VrArgu, reply *common.VrReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// do not process command if status is not NORMAL
	// and if i am not the primary in the current view
	if srv.status != NORMAL {
		return errors.New("status is INNORMAL")
	} else if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
		return errors.New("This is not Primary SRV")
	}

	command := args
	//append the command in its log
	srv.log = append(srv.log, command)
	go srv.resendPrepare(command, len(srv.log) - 1, srv.currentView, srv.commitIndex)

	//write -----> use
	op := args.Op
	switch op{
	case "DB.Login":
		//var temp *common.LogArgs
		temp,  _ := args.Argu.(common.LogArgs)
		var reply2 *common.LogReply
		srv.db.Login(&temp, reply2)
		reply.Reply = *reply2
	case "DB.Signup":
		temp, _ := args.Argu.(common.SignArgs)
		log.Print(temp)
		var reply2 common.SignReply
		srv.db.Signup(&temp, &reply2)
		log.Print("sign")
		reply.Reply = reply2
	case "DB.DelUser":
		temp, _ := args.Argu.(common.DelUserArgs)
		log.Print(temp)
		var reply2 common.DelUserReply
		srv.db.DelUser(&temp, &reply2)
		log.Print("delete user")
		reply.Reply = reply2
	case "DB.SendMsg":
		temp, _ := args.Argu.(common.SendMsgArgs)
		log.Print(temp)
		var reply2 common.SendMsgReply
		srv.db.SendMsg(&temp, &reply2)
		log.Print("snedmsg")
		reply.Reply = reply2
	case "DB.GetMsg":
		temp, _ := args.Argu.(common.GetMsgArgs)
		log.Print(temp)
		var reply2 common.GetMsgReply
		srv.db.GetMsg(&temp, &reply2)
		log.Print("getmsg")
		reply.Reply = reply2
	case "DB.LikeMsg":
		temp, _ := args.Argu.(common.LikeArgs)
		log.Print(temp)
		var reply2 common.LikeReply
		srv.db.LikeMsg(&temp, &reply2)
		log.Print("like")
		reply.Reply = reply2
	case "DB.UnLikeMsg":
		temp, _ := args.Argu.(common.UnLikeArgs)
		log.Print(temp)
		var reply2 common.UnLikeReply
		srv.db.UnLikeMsg(&temp, &reply2)
		log.Print("unlike")
		reply.Reply = reply2
	}



	log.Println("return")
//	log.Println(srv.IsCommitted(index))
	return nil
}

func (srv *PBServer) resendPrepare(command interface{}, index int, currentView int, commitIndex int) {
	replys := make([]*PrepareReply, len(srv.peers))
	for i, _ := range srv.peers {
		if i == srv.me {
			// replys[i] = &PrepareReply{
			// 	View:				srv.currentView,
			// 	Success:		true,
			// }
			continue
		}
		prepareArgs := &PrepareArgs{
			View:         currentView,         // the primary's current view
			PrimaryCommit:commitIndex, 			   // the primary's commitIndex
			Index:        index,        // the index position at which the log entry is to be replicated on backups
			Entry:        command,								 // the log entry to be replicated
		}
		replys[i] = &PrepareReply{}
		srv.sendPrepare(i, prepareArgs, replys[i])
	}
	//Recovery

	//change committedindex
	suc_num := 0
	for i, reply := range replys {
		if i == srv.me {
			continue
		}
		//log.Println(reply)
		if reply.Success == true {
			suc_num += 1
		}
	}
	//log.Printf("[%d] successful number", suc_num)

	if suc_num >= (len(srv.peers) - 1) / 2 {
		for {
			srv.mu.Lock()
			if srv.commitIndex + 1 == index {
				srv.commitIndex += 1
				srv.mu.Unlock()
				break
			}
			srv.mu.Unlock()
		}
	} else {
			//if not committed, resend prepare until index commmit
			go srv.resendPrepare(command, index, currentView, commitIndex)
	}


}

func (srv *PBServer) sendPrepare(server int, args *PrepareArgs, reply *PrepareReply) bool {
	ok := srv.peers[server].Call("PBServer.Prepare", args, reply)
	return ok==nil
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args *PrepareArgs, reply *PrepareReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.currentView == args.View && len(srv.log) == args.Index {
		srv.log = append(srv.log, args.Entry)
		reply.View = srv.currentView
		reply.Success = true
		srv.commitIndex = args.PrimaryCommit
		srv.lastNormalView = srv.currentView
		return
	}	else if srv.currentView == args.View && len(srv.log) > args.Index {
		//log's capaticy is bigger, so should return true, but not append,and not recover
		reply.View = args.View
		reply.Success = true
		return
	} else {
		//reply.View = srv.currentView
		reply.Success = false
		//no way to recover
		if(srv.currentView > args.View) {
			return
		}

		//deal with other prepare, no improve, don't know why.
		//go srv.sendPrepare(srv.me, args, reply)

		srv.status = RECOVERING
		go func() {
			recoveryArgs := &RecoveryArgs {
				View:			args.View,
				Server:		srv.me,
			}
			recoveryReply := &RecoveryReply{}
			ok := (srv.peers[GetPrimary(args.View, len(srv.peers))].Call("PBServer.Recovery", recoveryArgs, recoveryReply)==nil)
			if ok == true && recoveryReply.Success == true {
				srv.currentView = recoveryReply.View
				srv.lastNormalView = recoveryReply.View
				srv.commitIndex = recoveryReply.PrimaryCommit
				srv.log = recoveryReply.Entries
				srv.status = NORMAL
			}
		}()
	}
	// outdate := srv.currentView < args.View
	// log.Printf("if the date outdated",outdate)
	// if outdate {
	//
	// }


}

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(args *RecoveryArgs, reply *RecoveryReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.status == NORMAL {
		reply.View = srv.currentView
		reply.Entries = srv.log
		reply.PrimaryCommit = srv.commitIndex
		reply.Success = true
	} else {
		reply.Success = false
	}

}


func (srv *PBServer) PromptViewChange(newView int) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	newPrimary := GetPrimary(newView, len(srv.peers))

	if newPrimary != srv.me { //only primary of newView should do view change
		return
	} else if newView <= srv.currentView {
		return
	}
	vcArgs := &ViewChangeArgs{
		View: newView,
	}
	vcReplyChan := make(chan *ViewChangeReply, len(srv.peers))
	// send ViewChange to all servers including myself
	for i := 0; i < len(srv.peers); i++ {
		go func(server int) {
			var reply ViewChangeReply
			ok := srv.peers[server].Call("PBServer.ViewChange", vcArgs, &reply)==nil
			// fmt.Printf("node-%d (nReplies %d) received reply ok=%v reply=%v\n", srv.me, nReplies, ok, r.reply)
			if ok {
				vcReplyChan <- &reply
			} else {
				vcReplyChan <- nil
			}
		}(i)
	}

	// wait to receive ViewChange replies
	// if view change succeeds, send StartView RPC
	go func() {
		var successReplies []*ViewChangeReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil && r.Success {
				successReplies = append(successReplies, r)
			}
			if nReplies == len(srv.peers) || len(successReplies) == majority {
				break
			}
		}
		ok, log := srv.determineNewViewLog(successReplies)
		if !ok {
			return
		}
		svArgs := &StartViewArgs{
			View: vcArgs.View,
			Log:  log,
		}
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			var reply StartViewReply
			go func(server int) {
				// fmt.Printf("node-%d sending StartView v=%d to node-%d\n", srv.me, svArgs.View, server)
				srv.peers[server].Call("PBServer.StartView", svArgs, &reply)
			}(i)
		}
	}()
}


func (srv *PBServer) determineNewViewLog(successReplies []*ViewChangeReply) (
	ok bool, newViewLog []interface{}) {
		// Your code here

	max := 0
	if len(successReplies) >= len(srv.peers) / 2 {
		//decide new one
		ok = true
		//picks the log whose lastest normal view number is the largest.
		for i := 0; i < len(successReplies); i++ {
			if successReplies[i].LastNormalView > max {
				max = successReplies[i].LastNormalView
				//determined the log for the new-view
				newViewLog = successReplies[i].Log
				continue
			}
			if successReplies[i].LastNormalView == max {
				//more than one such logs, it picks the longest log among those
				if (len(successReplies[i].Log) > len(newViewLog)) {
					newViewLog = successReplies[i].Log
				}
			}
		}
	} else {
		ok = false
	}
	return ok, newViewLog
}

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer) ViewChange(args *ViewChangeArgs, reply *ViewChangeReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	// checks that the view number included in the message is indeed
	//larger than what it thinks the current view number is.
	if args.View > srv.currentView {
		srv.status = VIEWCHANGE
		reply.LastNormalView = srv.lastNormalView
		reply.Log = srv.log
		reply.Success = true
	} else {
		reply.Success = false
	}

}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(args *StartViewArgs, reply *StartViewReply) error {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	//the server must again check current view is no bigger than in the RPC message,
	//which would mean that there's been no concurrent view-change for a larger view.
	if srv.currentView > args.View {
		return nil
	} else {
		//sets the new-view as indicated in the message
		//changes its status to be NORMAL
		srv.log = args.Log
		srv.currentView = args.View
		//must set lastnormalview, other wise will return error exit status 1
		srv.lastNormalView = args.View
		srv.status = NORMAL
		srv.commitIndex = len(srv.log) - 1
	}
	return nil
}

// type GetServerNumberArgs struct{
//
// }
//
// type GetServerNumberReply struct{
// 	Number int
// }
//
// func (srv *PBServer) GetServerNumber(args *GetServerNumberArgs, reply *GetServerNumberReply) error{
// 	reply.Number = srv.me
// 	// log.Println("haha",srv.me)
// 	return nil
// }
//
//
// func main(){
// 	clients := make([]*rpc.Client, 1)
// 	// srv_num := 3
// 	// ports := []string{":8082",":8083",":8084"}
//
//
//
//
//
// 	port := flag.String("port", ":8080", "http listen port")
// 	num := flag.Int("num", 777, "client's number")
// 	flag.Parse()
//
// 	if(*port == ":8080" || *num == 777) {
// 		log.Print("! error: Please enter the parameter of port & num(client)")
// 		return
// 	}
// 	log.Print("port", *port)
//
// 	peer := Make(clients, *num, 0)
// 	server := rpc.NewServer()
// 	server.Register(peer)
// 	l,listenError := net.Listen("tcp", *port)
// 	if(listenError!=nil){
// 		log.Println(listenError)
// 	}
// 	go server.Accept(l)
//
// 	client, err := rpc.Dial("tcp", *port)
// 	clients[*num] = client
// 	if(err!=nil){
// 		log.Println(err)
// 	}
// 	log.Println(client==nil)
//
// 	signalChan := make(chan os.Signal, 1)
// 	cleanupDone := make(chan bool)
// 	signal.Notify(signalChan, os.Interrupt)
// 	go func() {
//     for _ = range signalChan {
//         log.Println("\nReceived an interrupt, stopping services...\n")
//         cleanupDone <- true
//     }
// 	}()
// 	<-cleanupDone
//
// }
