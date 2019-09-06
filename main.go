package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv/autoload"
	pr "gitlab.com/go-pher/go-auth/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// Server is a struct for grpc server, and contains various DB objects
type Server struct {
	coll  *mongo.Collection
	redis *redis.Client
}

func main() {
	// Context with timeout
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a listener on TCP port 60061
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 60061))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Mongo DB connection
	clientMongo, err := ConnectMongo(ctx)
	if err != nil {
		fmt.Println("ERROR: ", err)
		panic(err)
	}
	defer clientMongo.Disconnect(ctx)
	coll := clientMongo.Database("dev").Collection("users")

	// Redis connection
	clientRedis, err := ConnectRedis()
	if err != nil {
		fmt.Println("ERROR: ", err)
		panic(err)
	}
	defer clientRedis.Close()

	// Intialize Server strcut with both DB objects
	s := Server{
		coll:  coll,
		redis: clientRedis,
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	pr.RegisterAuthServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		panic(err)
	}
}
