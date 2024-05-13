package mr

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

const TIMEOUT = 5

// type ServerStatus int
// const (
// 		OK ServerStatus = iota
// 		WAITING
// 		FINISHED
// )
type JobStatus int

const (
	JOBSTATUS_SUCCESS JobStatus = iota
	JOBSTATUS_FAILED
	JOBSTATUS_INIT
	JOBSTATUS_MAP
	JOBSTATUS_REDUCE
)

type Coordinator struct {
	// Your definitions here.

	jobs_worker    map[string]int
	workers_status map[int]int // last online time
	jobs_status    map[string]JobStatus
}

func get_content(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	defer file.Close()
	return string(content)
}

func (c *Coordinator) get_new_workerid() int {
	var new_workerid int
	for {
		new_workerid = rand.Intn(8192 * 8192)
		value, ok := c.workers_status[new_workerid]
		if !ok || value+3*TIMEOUT < time.Now().Second() {
			break
		}
	}
	return new_workerid
}

func (c *Coordinator) register(args *RegisterArgs, reply *RegisterReply) error {
	reply.workerid = c.get_new_workerid()
	return nil
}

func (c *Coordinator) online_refresh(args *RefreshArgs, reply *RefreshReply) error {
	c.workers_status[args.workerid] = time.Now().Second()
	return nil
}

func (c *Coordinator) get_work(args *GetWorkArgs, reply *GetWorkReply) error {
	finished := true
	workerid := args.workerid
	for file, jobStatus := range c.jobs_status {
		if jobStatus == JOBSTATUS_INIT || jobStatus == JOBSTATUS_FAILED {
			c.jobs_status[file] = JOBSTATUS_MAP
			c.workers_status[workerid] = time.Now().Second()
			c.jobs_worker[file] = workerid
			reply.status = SERVERSTATUS_OK
			reply.kv = 
			return nil
		} else if jobStatus == JOBSTATUS_MAP || jobStatus == JOBSTATUS_REDUCE {
			finished = false
		}
	}
	if finished {
		reply.status = FINISHED
	} else {
		reply.status = WAITING
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.
	for _, file := range files {
		c.jobs_status[file] = INIT
	}

	c.server()
	return &c
}
