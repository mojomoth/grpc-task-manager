package grpctaskmanager

import (
	"context"
	"fmt"
	"time"

	pb "github.com/dev-zipida-com/grpc-task-manager/protos/v1/task"
)

type TaskQueue struct {
	queue      Queue
	roundRobin int
	clients    []*pb.Client
	streams    []*pb.Notification_SubscribeServer
}

func (taskQueue *TaskQueue) initialize() {
	taskQueue.queue = Queue{}
	taskQueue.roundRobin = 0
	taskQueue.clients = make([]*pb.Client, 0)
	taskQueue.streams = make([]*pb.Notification_SubscribeServer, 0)
}

func (taskQueue *TaskQueue) destroy(client *pb.Client, stream *pb.Notification_SubscribeServer) {
	taskQueue.changeRoundRobin(stream)
	taskQueue.removeClient(client)
	taskQueue.removeStream(stream)
}

func (taskQueue *TaskQueue) changeRoundRobin(stream *pb.Notification_SubscribeServer) {
	var index int
	for i, v := range taskQueue.streams {
		if v == stream {
			index = i
			break
		}
	}

	if taskQueue.roundRobin > index {
		taskQueue.roundRobin--
	}

	if len(taskQueue.streams)-1 == index {
		taskQueue.roundRobin = 0
	}
}

func (taskQueue *TaskQueue) removeClient(client *pb.Client) {
	var index int
	s := taskQueue.clients
	for i, v := range s {
		if v == client {
			index = i
			break
		}
	}

	sliced := append(s[:index], s[index+1:]...)
	taskQueue.clients = sliced
}

func (taskQueue *TaskQueue) removeStream(stream *pb.Notification_SubscribeServer) {
	var index int
	s := taskQueue.streams
	for i, v := range s {
		if v == stream {
			index = i
			break
		}
	}

	sliced := append(s[:index], s[index+1:]...)
	taskQueue.streams = sliced
}

func (taskQueue *TaskQueue) deliver() {
	if taskQueue.queue.IsEmpty() {
		return
	}

	if len(taskQueue.streams) == 0 {
		return
	}

	// send message
	message := taskQueue.queue.Dequeue()
	(*taskQueue.streams[taskQueue.roundRobin]).Send(&pb.SubscribeResponse{
		Task: &pb.Task{
			Message: fmt.Sprintf("%v", message),
			Time:    time.Now().UnixNano(),
		},
		IsOk: true,
	})

	// update index
	taskQueue.roundRobin++
	if taskQueue.roundRobin == len(taskQueue.streams) {
		taskQueue.roundRobin = 0
	}
}

type TaskManager struct {
	clients    map[string]*pb.Client
	streams    map[string]*pb.Notification_SubscribeServer
	taskQueues map[string]*TaskQueue
	pb.UnimplementedNotificationServer
}

func (taskManager *TaskManager) initialize() {
	taskManager.clients = make(map[string]*pb.Client)
	taskManager.streams = make(map[string]*pb.Notification_SubscribeServer)
	taskManager.taskQueues = make(map[string]*TaskQueue)
}

func (taskManager *TaskManager) destroy(taskId string, clientId string) {
	taskManager.taskQueues[taskId].destroy(taskManager.clients[clientId], taskManager.streams[clientId])
	delete(taskManager.clients, clientId)
	delete(taskManager.streams, clientId)

	if len(taskManager.taskQueues[taskId].queue) == 0 {
		delete(taskManager.taskQueues, taskId)
	}
}

func (taskManager *TaskManager) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	id := req.Client.Id
	taskManager.clients[id] = req.Client

	return &pb.ConnectResponse{
		IsOk: true,
	}, nil
}

func (taskManager *TaskManager) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	taskQueue := taskManager.checkExistAndCreate(req.TaskId)
	taskQueue.queue.Enqueue(req.Message)

	if len(taskQueue.streams) > 0 {
		taskQueue.deliver()
	}

	return &pb.PublishResponse{
		IsOk: true,
	}, nil
}

func (taskManager *TaskManager) Subscribe(req *pb.SubscribeRequest, stream pb.Notification_SubscribeServer) error {
	taskQueue := taskManager.checkExistAndCreate(req.TaskId)

	clientId := req.Client.Id
	taskManager.streams[clientId] = &stream

	// add
	taskQueue.clients = append(taskQueue.clients, req.Client)
	taskQueue.streams = append(taskQueue.streams, &stream)

	for {
		if !taskQueue.queue.IsEmpty() {
			taskQueue.deliver()
		}

		err := stream.Context().Err()
		if err != nil {
			break
		}
	}

	// destroy
	taskManager.destroy(req.TaskId, req.Client.Id)

	return nil
}

func (taskManager *TaskManager) checkExistAndCreate(taskId string) *TaskQueue {
	if _, exist := taskManager.taskQueues[taskId]; !exist {
		taskQueue := &TaskQueue{}
		taskQueue.initialize()
		taskManager.taskQueues[taskId] = taskQueue
	}

	return taskManager.taskQueues[taskId]
}
