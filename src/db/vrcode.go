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
	// "fmt"
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
	gob.Register(common.VrArgu{})
  gob.Register(common.VrReply{})
  gob.Register(common.SignReply{})
  gob.Register(common.SignArgs{})
  gob.Register(common.LogArgs{})
  gob.Register(common.LogReply{})
  gob.Register(common.DelUserArgs{})
  gob.Register(common.DelUserReply{})
  gob.Register(common.SendMsgArgs{})
  gob.Register(common.SendMsgReply{})
  gob.Register(common.GetMsgArgs{})
  gob.Register(common.GetMsgReply{})
  gob.Register(common.LikeArgs{})
  gob.Register(common.LikeReply{})
  gob.Register(common.UnLikeArgs{})
  gob.Register(common.UnLikeReply{})
  gob.Register(common.LikeListArgs{})
  gob.Register(common.LikeListReply{})
  gob.Register(common.IsLikeArgs{})
  gob.Register(common.IsLikeReply{})
  gob.Register(common.FollowUserArgs{})
  gob.Register(common.FollowUserReply{})
  gob.Register(common.FollowListArgs{})
  gob.Register(common.FollowListReply{})

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

func (srv *PBServer) DealPrimay(args common.DealPrimayArgs, reply *common.DealPrimayReply) error {
	if (srv.me == GetPrimary(srv.currentView, len(srv.peers)) && srv.status==NORMAL){
		reply.OK = true
	} else {
		reply.OK = false
	}
	return nil
}


func (srv *PBServer) commitDB(args common.VrArgu) interface{} {
	op := args.Op
	switch op{
	case "DB.Login":
		temp,  _ := args.Argu.(common.LogArgs)
		var reply2 common.LogReply
		srv.db.Login(&temp, &reply2)
		return reply2
	case "DB.Signup":
		temp, _ := args.Argu.(common.SignArgs)
		log.Print(temp)
		var reply2 common.SignReply
		srv.db.Signup(&temp, &reply2)
		log.Print("sign")
		return reply2
	case "DB.DelUser":
		temp, _ := args.Argu.(common.DelUserArgs)
		log.Print(temp)
		var reply2 common.DelUserReply
		srv.db.DelUser(&temp, &reply2)
		log.Print("delete user")
		return reply2
	case "DB.SendMsg":
		temp, _ := args.Argu.(common.SendMsgArgs)
		log.Print(temp)
		var reply2 common.SendMsgReply
		srv.db.SendMsg(&temp, &reply2)
		log.Print("snedmsg")
		return reply2
	case "DB.GetMsg":
		temp, _ := args.Argu.(common.GetMsgArgs)
		log.Print(temp)
		var reply2 common.GetMsgReply
		srv.db.GetMsg(&temp, &reply2)
		log.Print("getmsg")
		return reply2
	case "DB.LikeMsg":
		temp, _ := args.Argu.(common.LikeArgs)
		log.Print(temp)
		var reply2 common.LikeReply
		srv.db.LikeMsg(&temp, &reply2)
		log.Print("like")
		return reply2
	case "DB.UnLikeMsg":
		temp, _ := args.Argu.(common.UnLikeArgs)
		log.Print(temp)
		var reply2 common.UnLikeReply
		srv.db.UnLikeMsg(&temp, &reply2)
		log.Print("unlike")
		return reply2
	case "DB.IsLike":
		temp, _ := args.Argu.(common.IsLikeArgs)
		log.Print(temp)
		var reply2 common.IsLikeReply
		srv.db.IsLike(&temp, &reply2)
		log.Print("IsLike")
		return reply2
	case "DB.FollowList":
		temp, _ := args.Argu.(common.FollowListArgs)
		log.Print(temp)
		var reply2 common.FollowListReply
		srv.db.FollowList(&temp, &reply2)
		log.Print("FollowList")
		return reply2
	case "DB.FollowUser":
		temp, _ := args.Argu.(common.FollowUserArgs)
		log.Print(temp)
		var reply2 common.FollowUserReply
		srv.db.FollowUser(&temp, &reply2)
		log.Print("FollowUser")
		return reply2
	case "DB.LikeList":
		temp, _ := args.Argu.(common.LikeListArgs)
		log.Print(temp)
		var reply2 common.LikeListReply
		srv.db.LikeList(&temp, &reply2)
		log.Print("LikeList:", reply2)
		return reply2
	}

	return 0
}

func (srv *PBServer) Start(args common.VrArgu, reply *common.VrReply) error {
	log.Println("Start run")
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
	reply.Reply = srv.resendPrepare(command, len(srv.log) - 1, srv.currentView, srv.commitIndex)
	log.Println("Start return")
	return nil
}

func (srv *PBServer) resendPrepare(command interface{}, index int, currentView int, commitIndex int) interface{}{
	// log.Println("start resendPrepare")
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
	// log.Println("suc_num:", suc_num)
	// log.Println("len(srv.peers)", len(srv.peers))
	if suc_num >= (len(srv.peers) - 1) / 2 {
		for {
			if srv.commitIndex + 1 == index {
				srv.commitIndex += 1
				argu, _ := command.(common.VrArgu)
				return srv.commitDB(argu)
			}
		}
	} else {
			//if not committed, resend prepare until index commmit
			// log.Println("resendPrepare call itself")
			return srv.resendPrepare(command, index, currentView, commitIndex)
	}

	return 1
}

func (srv *PBServer) sendPrepare(server int, args *PrepareArgs, reply *PrepareReply) bool {
	// log.Println("sendPrepare:", srv.peers[server])
	// log.Println("sendPrepare num:", len(srv.peers))
	prepareErr := srv.peers[server].Call("PBServer.Prepare", args, reply)
	// log.Println("prepareErr:",prepareErr)
	return prepareErr==nil
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args *PrepareArgs, reply *PrepareReply) error {
	// Your code here
	log.Println("start prepare from:", srv.me)
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.currentView == args.View && len(srv.log) == args.Index {
		log.Println("prepare1 from:", srv.me)
		log.Println("Prepare1")
		srv.log = append(srv.log, args.Entry)
		reply.View = srv.currentView
		reply.Success = true
		// log.Println("srv.commitIndex", srv.commitIndex)
		// log.Println("args.PrimaryCommit", args.PrimaryCommit)
		for i:=srv.commitIndex+1; i<=args.PrimaryCommit; i++{
			// log.Println("Prepare commit index:", i)
			argu, _ := srv.log[i].(common.VrArgu)
			srv.commitDB(argu)
		}
		srv.commitIndex = args.PrimaryCommit

		// srv.commitDB()
		srv.lastNormalView = srv.currentView
		return nil
	}	else if srv.currentView == args.View && len(srv.log) > args.Index {
		log.Println("prepare2 from:", srv.me)
		//log's capaticy is bigger, so should return true, but not append,and not recover
		log.Println("Prepare2")
		reply.View = args.View
		reply.Success = true
		return nil
	} else {
		log.Println("prepare3 from:", srv.me)
		//reply.View = srv.currentView
		reply.Success = false
		//no way to recover
		if(srv.currentView > args.View) {
			return nil
		}

		//deal with other prepare, no improve, don't know why.
		//go srv.sendPrepare(srv.me, args, reply)
		log.Println("Prepare4")
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
		return nil
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

func (srv *PBServer) VrViewChange(args *common.VrViewChangeArgu, reply *common.VrViewChangeReply) error{
	srv.PromptViewChange(args.View)
	return nil
}

func (srv *PBServer) PromptViewChange(newView int) {
	log.Println("PromptViewChange from:", srv.me)
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
		log.Println("middle1 PromptViewChange from:", srv.me)
		go func(server int) {
			var reply ViewChangeReply
			vcError := srv.peers[server].Call("PBServer.ViewChange", vcArgs, &reply)
			if(vcError!=nil){
				log.Println("vcError",vcError)
			}
			ok := vcError==nil
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
			// log.Println("middle2 PromptViewChange from:", srv.me)
			// log.Println("middle2.1 PromptViewChange from:", r)
			nReplies++
			if r != nil && r.Success {
				// log.Println("middle3 PromptViewChange from:", srv.me)
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

		// fmt.Println("middle4 PromptViewChange from:", srv.me)
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			var reply StartViewReply
			go func(server int) {
				// fmt.Println("middle4 PromptViewChange to:", server)
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
func (srv *PBServer) ViewChange(args *ViewChangeArgs, reply *ViewChangeReply) error {
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
	return nil
}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(args *StartViewArgs, reply *StartViewReply) error {
	// log.Println("StartView from:", srv.me)
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	//the server must again check current view is no bigger than in the RPC message,
	//which would mean that there's been no concurrent view-change for a larger view.
	if srv.currentView > args.View {
		return nil
	} else {
		// log.Println("StartView2 from:", srv.me)
		//sets the new-view as indicated in the message
		//changes its status to be NORMAL
		srv.log = args.Log
		srv.currentView = args.View
		//must set lastnormalview, other wise will return error exit status 1
		srv.lastNormalView = args.View
		srv.status = NORMAL
		srv.commitIndex = len(srv.log) - 1
		// log.Println("StartView3 from:", srv.currentView)
	}
	return nil
}
