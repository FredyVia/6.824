package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"time"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func call_register() int {
	args := RegisterArgs{}
	reply := RegisterReply{}

	ok := call("Coordinator.register", &args, &reply)
	if !ok {
		fmt.Printf("register call failed!\n")
	}
	return reply.workerid
}

func call_refresh(workerid int) {
	args := RefreshArgs{}

	// fill in the argument(s).
	args.workerid = workerid

	// declare a reply structure.

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.online_refresh", &args, nil)
	if !ok {
		fmt.Printf("online refresh call failed!\n")
	}
}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	workerid := call_register()
	// Your worker implementation here.
	ticker := time.NewTicker(TIMEOUT)
	stopChan := make(chan struct{}) // 创建用于发送停止信号的通道

	go func() {
		for {
			select {
			case <-ticker.C:
				call_refresh(workerid)
			case <-stopChan:
				fmt.Println("Ticker stopped")
				ticker.Stop() // 停止ticker
				return        // 退出goroutine
			}
		}
	}()
	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	for {
		args := GetWorkArgs{}
		args.workerid = workerid
		// declare a reply structure.
		reply := GetWorkReply{}

		ok := call("Coordinator.get_work", nil, &reply)
		if !ok {
			fmt.Printf("GetWork call failed!\n")
			break
		}

		if reply.status == FINISHED {
			break
		} else if reply.status == WAITING {
			time.Sleep(5000)
			continue
		}

		mapf(reply.key)
		reply.key
	}
	stopChan <- struct{}{} // 发送停止信号
	fmt.Println("Stop signal sent")
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
