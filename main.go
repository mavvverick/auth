package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	pr "gitlab.com/go-pher/go-auth/proto"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Server is a struct for grpc server, and contains various DB objects
type Server struct {
	db    *gorm.DB
	redis *redis.Client
}

func main() {
	// Context with timeout
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a listener on TCP port 60061
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 60061))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	time.Sleep(5 * time.Second)
	// Print the IP of the POD
	// cmd, err := exec.Command("curl", "ifconfig.co").Output()
	// fmt.Println("OUTPUT IP---", string(cmd))

	// MySQL connection
	clientSQL, err := ConnectSQL()
	if err != nil {
		fmt.Println("ERROR: ", err)
		panic(err)
	}
	defer clientSQL.Close()

	// Redis connection
	clientRedis, err := ConnectRedis()
	if err != nil {
		fmt.Println("ERROR: ", err)
		panic(err)
	}
	defer clientRedis.Close()

	// Intialize Server strcut with both DB objects
	s := Server{
		db:    clientSQL,
		redis: clientRedis,
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	pr.RegisterAuthServer(grpcServer, &s)
	health.RegisterHealthServer(grpcServer, s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		panic(err)
	}
}
