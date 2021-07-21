package grpctaskmanager

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/dev-zipida-com/grpc-task-manager/protos/v1/task"
)

const PORT = "9000"
const MESSAGE_SIZE = 100 * 1024 * 1024

func Run() {
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// logic
	taskManager := &TaskManager{}
	taskManager.initialize()

	// server
	options := []grpc.ServerOption{}
	options = append(options, grpc.MaxMsgSize(MESSAGE_SIZE))
	options = append(options, grpc.MaxRecvMsgSize(MESSAGE_SIZE))
	// options = append(options, grpc.StatsHandler(&Handler{}))
	grpcServer := grpc.NewServer(options...)

	// register
	pb.RegisterNotificationServer(grpcServer, taskManager)
	go TestStatus(taskManager)

	// serve
	log.Printf("start server on %s port", PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// 1초에 한번씩 체크하는 함수 만들어보기
func TestStatus(taskManager *TaskManager) {
	for {
		for k, q := range taskManager.taskQueues {
			log.Printf("task-key: %s, client-length: %d, ready-index : %d, queue: %v", k, len(q.streams), q.roundRobin, q.queue)
		}

		time.Sleep(time.Second)
	}
}
