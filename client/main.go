package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/proto1"
	"math/rand"
	"time"

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

func NormalAbnormalHeartBeat(c proto1.HeartBeatServiceClient) {
	stream, err := c.NormalAbnormalHeartBeat(context.Background())
	handleError(err)
	for t := 0; t < 10; t++ {
		newNARequst := &proto1.NormalAbnormalHeartBeatRequest{
			Bpm: generateBPM(),
		}
		stream.Send(newNARequst)
		fmt.Printf("Sent %v\n", newNARequst)
	}
	stream.CloseSend()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		handleError(err)
		time.Sleep(1 * time.Second)
		fmt.Printf("Received %v\n", msg)
	}

}

func HeartBeatHistory(c proto1.HeartBeatServiceClient) {
	newHistoryRequest := proto1.HeartBeatHistoryRequest{
		Username: "hemanthghs",
	}
	res_stream, err := c.HeartBeatHistory(context.Background(), &newHistoryRequest)
	handleError(err)
	for {
		msg, err := res_stream.Recv()
		handleError(err)
		if err == io.EOF {
			break
		}
		fmt.Println(msg)
		// time.Sleep(1 * time.Second)
	}
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
		fmt.Println("Resuest send ", newRequest)
		stream.Send(newRequest)
	}

	stream.CloseAndRecv()
	// handleError(err)
	// fmt.Println(resp)

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
	// UserHeartBeat(c)
	// LiveHeartBeat(c)
	// HeartBeatHistory(c)
	NormalAbnormalHeartBeat(c)
}
