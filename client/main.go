package main

import (
	"context"
	"fmt"
	"log"
	"main/proto1"
	"math/rand"

	"google.golang.org/grpc"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateBPM() int32 {
	bpm := rand.Intn(100)
	return int32(bpm)
}

func LiveHeartBeat(c proto1.HeartBeatServiceClient) {
	stream, err := c.LiveHeartBeat(context.Background())
	handleError(err)
	for t := 0; t < 10; t++ {
		newRequest := &proto1.LiveHeartBeatRequest{
			Heartbeat: &proto1.HeartBeat{
				Bpm:      generateBPM(),
				Username: "hemanthghs",
			},
		}
		stream.Send(newRequest)
	}

	resp, err := stream.CloseAndRecv()
	handleError(err)
	fmt.Println(resp)

}

func UserHeartBeat(c proto1.HeartBeatServiceClient) {
	heartbeatRequest := proto1.HeartBeatRequest{
		Heartbeat: &proto1.HeartBeat{
			Bpm:      65,
			Username: "hemanth",
		},
	}
	c.UserHeartBeat(context.Background(), &heartbeatRequest)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	handleError(err)
	fmt.Println("Client started")
	defer conn.Close()

	c := proto1.NewHeartBeatServiceClient(conn)
	UserHeartBeat(c)
}
