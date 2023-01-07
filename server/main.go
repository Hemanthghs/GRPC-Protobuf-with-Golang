package main

import (
	"context"
	"fmt"
	"log"
	"main/proto1"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	proto1.UnimplementedHeartBeatServiceServer
}

type heart_item struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Bpm      int32              `bson:"bpm"`
	Username string             `bson:"username"`
}

func pushUserToDb(ctx context.Context, item heart_item) primitive.ObjectID {
	res, err := collection.InsertOne(ctx, item)
	handleError(err)

	return res.InsertedID.(primitive.ObjectID)
}

func (*server) UserHeartBeat(ctx context.Context, req *proto1.HeartBeatRequest) (*proto1.HeartBeatResponse, error) {
	fmt.Println(req)
	bpm := req.GetHeartbeat().GetBpm()
	username := req.GetHeartbeat().GetUsername()
	newHeartItem := heart_item{
		Bpm:      bpm,
		Username: username,
	}
	docid := pushUserToDb(ctx, newHeartItem)
	result := fmt.Sprintf("User HeartBeat is %v, newly created docid is %v", bpm, docid)
	heartBeatResponse := proto1.HeartBeatResponse{
		Result: result,
	}
	return &heartBeatResponse, nil
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var collection *mongo.Collection

func main() {

	godotenv.Load(".env")

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	proto1.RegisterHeartBeatServiceServer(s, &server{})

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	handleError(err)

	go func() {
		if err := s.Serve(lis); err != nil {
			handleError(err)
		}
	}()
	mongo_uri := goDotEnvVariable("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	handleError(err)
	fmt.Println("Mongo connected")
	err = client.Connect(context.TODO())
	handleError(err)

	collection = client.Database("heartbeat").Collection("heartbeat")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	fmt.Println("closing mongo connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		handleError(err)
	}

	s.Stop()
}
