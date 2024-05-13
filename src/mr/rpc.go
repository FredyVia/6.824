package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

// Add your RPC definitions here.
type RegisterArgs struct {
	workerid int
}

type RegisterReply struct {
	workerid int
}

type RefreshArgs struct {
	workerid int
}

type RefreshReply struct {
}

type ServerStatus int
type WorkType int

const (
	SERVERSTATUS_OK ServerStatus = iota
	SERVERSTATUS_WAITING
	SERVERSTATUS_FINISHED
)

const (
	WORKTYPE_MAP WorkType = iota
	WORKTYPE_REDUCE
)

type GetWorkArgs struct {
	workerid int
}

type GetWorkReply struct {
	status   ServerStatus
	worktype WorkType
	kv       KeyValue
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
