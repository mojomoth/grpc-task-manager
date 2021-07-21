package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/dev-zipida-com/grpc-task-manager/protos/v1/task"
)

const TEST_ADDRESS = "localhost:9000"

func main() {
	// connect
	conn, err := grpc.Dial(TEST_ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// notification
	cli := pb.NewNotificationClient(conn)

	// publish
	cli.Publish(context.Background(), &pb.PublishRequest{
		TaskId:  "test",
		Message: "{\n\"foo\": \"bar\"\n}",
	})
}
