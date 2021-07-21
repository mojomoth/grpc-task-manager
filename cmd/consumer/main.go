package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "github.com/dev-zipida-com/grpc-task-manager/protos/v1/task"
)

const TEST_ADDRESS = "localhost:9000"

func main() {
	uuid := uuid.New().String()
	client := &pb.Client{
		Id:   uuid,
		Type: "test",
		Name: "name-" + uuid,
		Time: time.Now().UnixNano(),
	}

	// connect
	conn, err := grpc.Dial(TEST_ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// notification
	cli := pb.NewNotificationClient(conn)

	// connect
	cli.Connect(context.Background(), &pb.ConnectRequest{Client: client})

	// subscribe
	stream, _ := cli.Subscribe(context.Background(), &pb.SubscribeRequest{
		Client: client,
		TaskId: "test",
	})

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// no more stream to listen
			break
		}
		if err != nil {
			// some error occured
			log.Fatalf("%v", err)
		}

		// recv message
		log.Printf("receive message: %s (time: %d)", res.Task.Message, res.Task.Time/1e6)
	}
}
