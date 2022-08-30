package service

import (
	"dataProcess/constants"
	"dataProcess/constants/statuses"
	"dataProcess/service/local/impl"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"dataProcess/tcp_socket/socket"
	"dataProcess/tools/utils/log"
)

var cancelTask = make(chan string, 1000)
var reply socket.TaskReply

func InitServerApp(network, host, port string) {
	svr := socket.Server{NewClient: newClient}
	svr.Start(network, host+constants.Colon+port)
	svr.RegisterMethod("receiveProcess", receiveProcess)
	fmt.Printf("ServerApp Run http://%s:%s/ \r\n", host, port)
	fmt.Printf("Enter Control + C Shutdown Server \r\n")
	for {
		select {
		case id := <-cancelTask:
			call := socket.Call{Method: "cancel", Args: &socket.CancelArgs{Cause: "test"}, Reply: &socket.CancelReply{}}
			if s, ok := svr.GetClient(id); ok {
				log.Debug("%v", s.Call(&call))
			}
			log.Debug("%v", call.Reply.(*socket.CancelReply).Result)
		default:
			time.Sleep(time.Second)
		}
	}
}

func receiveProcess(id string, p *socket.Packet) []byte {
	log.Debug("receiveProcess: %s,%s", id, p.Msg.Params.String())

	var proReply socket.ProcessReply
	err := json.Unmarshal([]byte(p.Msg.Params.String()), &proReply)
	if err != nil {
		log.Debug("err:", err.Error())
		os.Exit(1)
	}

	taskId, err := utils.StringToint32(proReply.TaskId)
	if err != nil {
		log.Debug("err:", err.Error())
		os.Exit(1)
	}

	impl.CreateTaskManageImpl().UpdateTaskStatusEnter(taskId, proReply.Status, proReply.Stage, proReply.Exception, proReply.EndTime)
	return cancelCmdDes(id, p, &proReply)
}

func newClient(id string) {
	log.Debug(id)
}

func cancelCmdDes(id string, p *socket.Packet, proReply *socket.ProcessReply) []byte {
	reply.TaskId = proReply.TaskId
	reply.ErrMsg = proReply.Exception

	switch proReply.Status {
	case statuses.Error, statuses.Failed:
		reply.Result = false
	case statuses.Terminal:
		reply.Result = false
		cancelTask <- id
	default:
		reply.Result = true
	}
	return socket.MarshalRespond(p.Msg.Id, &reply, p.MessageBuf()) // 发回的消息格式
}
